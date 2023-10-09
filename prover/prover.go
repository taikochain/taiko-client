package prover

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/taikoxyz/taiko-client/bindings"
	"github.com/taikoxyz/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-client/metrics"
	eventIterator "github.com/taikoxyz/taiko-client/pkg/chain_iterator/event_iterator"
	"github.com/taikoxyz/taiko-client/pkg/rpc"
	capacity "github.com/taikoxyz/taiko-client/prover/capacity_manager"
	proofProducer "github.com/taikoxyz/taiko-client/prover/proof_producer"
	proofSubmitter "github.com/taikoxyz/taiko-client/prover/proof_submitter"
	"github.com/taikoxyz/taiko-client/prover/server"
	"github.com/urfave/cli/v2"
)

var (
	errNoCapacity   = errors.New("no prover capacity available")
	errTierNotFound = errors.New("tier not found")
)

// Prover keep trying to prove new proposed blocks valid/invalid.
type Prover struct {
	// Configurations
	cfg           *Config
	proverAddress common.Address

	// Clients
	rpc *rpc.Client

	// Prover Server
	srv *server.ProverServer

	// Contract configurations
	protocolConfigs *bindings.TaikoDataConfig

	// States
	latestVerifiedL1Height uint64
	lastHandledBlockID     uint64
	genesisHeightL1        uint64
	l1Current              *types.Header
	reorgDetectedFlag      bool
	tiers                  []*rpc.TierProviderTierWithID

	// Proof submitters
	proofSubmitters []proofSubmitter.Submitter

	// Subscriptions
	blockProposedCh      chan *bindings.TaikoL1ClientBlockProposed
	blockProposedSub     event.Subscription
	transitionProvedCh   chan *bindings.TaikoL1ClientTransitionProved
	transitionProvedSub  event.Subscription
	blockVerifiedCh      chan *bindings.TaikoL1ClientBlockVerified
	blockVerifiedSub     event.Subscription
	proofWindowExpiredCh chan *bindings.TaikoL1ClientBlockProposed
	proveNotify          chan struct{}

	// Proof related
	proofGenerationCh chan *proofProducer.ProofWithHeader

	// Concurrency guards
	proposeConcurrencyGuard     chan struct{}
	submitProofConcurrencyGuard chan struct{}

	// capacity-related configs
	capacityManager *capacity.CapacityManager

	ctx context.Context
	wg  sync.WaitGroup
}

// InitFromCli initializes the given prover instance based on the command line flags.
func (p *Prover) InitFromCli(ctx context.Context, c *cli.Context) error {
	cfg, err := NewConfigFromCliContext(c)
	if err != nil {
		return err
	}

	return InitFromConfig(ctx, p, cfg)
}

// InitFromConfig initializes the prover instance based on the given configurations.
func InitFromConfig(ctx context.Context, p *Prover, cfg *Config) (err error) {
	p.cfg = cfg
	p.ctx = ctx
	p.capacityManager = capacity.New(cfg.Capacity)

	// Clients
	if p.rpc, err = rpc.NewClient(p.ctx, &rpc.ClientConfig{
		L1Endpoint:        cfg.L1WsEndpoint,
		L2Endpoint:        cfg.L2WsEndpoint,
		TaikoL1Address:    cfg.TaikoL1Address,
		TaikoL2Address:    cfg.TaikoL2Address,
		TaikoTokenAddress: cfg.TaikoTokenAddress,
		RetryInterval:     cfg.BackOffRetryInterval,
		Timeout:           cfg.RPCTimeout,
		BackOffMaxRetrys:  new(big.Int).SetUint64(p.cfg.BackOffMaxRetrys),
	}); err != nil {
		return err
	}

	// Configs
	protocolConfigs, err := p.rpc.TaikoL1.GetConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get protocol configs: %w", err)
	}
	p.protocolConfigs = &protocolConfigs

	log.Info("Protocol configs", "configs", p.protocolConfigs)

	p.proverAddress = crypto.PubkeyToAddress(p.cfg.L1ProverPrivKey.PublicKey)

	chBufferSize := p.protocolConfigs.BlockMaxProposals
	p.blockProposedCh = make(chan *bindings.TaikoL1ClientBlockProposed, chBufferSize)
	p.blockVerifiedCh = make(chan *bindings.TaikoL1ClientBlockVerified, chBufferSize)
	p.transitionProvedCh = make(chan *bindings.TaikoL1ClientTransitionProved, chBufferSize)
	p.proofGenerationCh = make(chan *proofProducer.ProofWithHeader, chBufferSize)
	p.proofWindowExpiredCh = make(chan *bindings.TaikoL1ClientBlockProposed, chBufferSize)
	p.proveNotify = make(chan struct{}, 1)
	if err := p.initL1Current(cfg.StartingBlockID); err != nil {
		return fmt.Errorf("initialize L1 current cursor error: %w", err)
	}

	// Concurrency guards
	p.proposeConcurrencyGuard = make(chan struct{}, cfg.MaxConcurrentProvingJobs)
	p.submitProofConcurrencyGuard = make(chan struct{}, cfg.MaxConcurrentProvingJobs)

	// Protocol proof tiers
	if p.tiers, err = p.rpc.GetTiers(ctx); err != nil {
		return err
	}

	// Proof submitters
	for _, tier := range p.tiers {
		var (
			producer  proofProducer.ProofProducer
			submitter proofSubmitter.Submitter
		)
		switch tier.ID {
		case encoding.TierOptimisticID:
			producer = &proofProducer.OptimisticProofProducer{DummyProofProducer: new(proofProducer.DummyProofProducer)}
		case encoding.TierSgxID:
			producer = &proofProducer.SGXProofProducer{DummyProofProducer: new(proofProducer.DummyProofProducer)}
		case encoding.TierPseZkevmID:
			zkEvmRpcdProducer, err := proofProducer.NewZkevmRpcdProducer(
				cfg.ZKEvmRpcdEndpoint,
				cfg.ZkEvmRpcdParamsPath,
				cfg.L1HttpEndpoint,
				cfg.L2HttpEndpoint,
				true,
				p.protocolConfigs,
			)
			if err != nil {
				return err
			}
			if p.cfg.Dummy {
				zkEvmRpcdProducer.DummyProofProducer = new(proofProducer.DummyProofProducer)
			}
			producer = zkEvmRpcdProducer
		case encoding.TierGuardianID:
			producer = &proofProducer.GuardianProofProducer{DummyProofProducer: new(proofProducer.DummyProofProducer)}
		}

		if submitter, err = proofSubmitter.New(
			p.rpc,
			producer,
			p.proofGenerationCh,
			p.cfg.TaikoL2Address,
			p.cfg.L1ProverPrivKey,
			p.cfg.Graffiti,
			p.cfg.ProofSubmissionMaxRetry,
			p.cfg.BackOffRetryInterval,
			p.cfg.WaitReceiptTimeout,
			p.cfg.ProveBlockGasLimit,
			p.cfg.ProveBlockTxReplacementMultiplier,
			p.cfg.ProveBlockMaxTxGasTipCap,
		); err != nil {
			return err
		}

		p.proofSubmitters = append(p.proofSubmitters, submitter)
	}

	// Prover server
	proverServerOpts := &server.NewProverServerOpts{
		ProverPrivateKey:     p.cfg.L1ProverPrivKey,
		MinOptimisticTierFee: p.cfg.MinOptimisticTierFee,
		MinSgxTierFee:        p.cfg.MinSgxTierFee,
		MinPseZkevmTierFee:   p.cfg.MinPseZkevmTierFee,
		MaxExpiry:            p.cfg.MaxExpiry,
		CapacityManager:      p.capacityManager,
		TaikoL1Address:       p.cfg.TaikoL1Address,
		Rpc:                  p.rpc,
		LivenessBond:         protocolConfigs.LivenessBond,
		IsGuardian:           p.cfg.GuardianProver,
	}
	if p.cfg.GuardianProver {
		proverServerOpts.ProverPrivateKey = p.cfg.GuardianProverPrivateKey
	}
	if p.srv, err = server.New(proverServerOpts); err != nil {
		return err
	}

	return nil
}

// Start starts the main loop of the L2 block prover.
func (p *Prover) Start() error {
	p.wg.Add(1)
	p.initSubscription()
	go func() {
		if err := p.srv.Start(fmt.Sprintf(":%v", p.cfg.HTTPServerPort)); !errors.Is(err, http.ErrServerClosed) {
			log.Crit("Failed to start http server", "error", err)
		}
	}()
	go p.eventLoop()

	return nil
}

// eventLoop starts the main loop of Taiko prover.
func (p *Prover) eventLoop() {
	defer func() {
		p.wg.Done()
	}()

	// reqProving requests performing a proving operation, won't block
	// if we are already proving.
	reqProving := func() {
		select {
		case p.proveNotify <- struct{}{}:
		default:
		}
	}

	// If there is too many (TaikoData.Config.blockMaxProposals) pending blocks in TaikoL1 contract, there will be no new
	// BlockProposed temporarily, so except the BlockProposed subscription, we need another trigger to start
	// fetching the proposed blocks.
	forceProvingTicker := time.NewTicker(15 * time.Second)
	defer forceProvingTicker.Stop()

	// Call reqProving() right away to catch up with the latest state.
	reqProving()

	for {
		select {
		case <-p.ctx.Done():
			return
		case proofWithHeader := <-p.proofGenerationCh:
			p.submitProofOp(p.ctx, proofWithHeader)
		case <-p.proveNotify:
			if err := p.proveOp(); err != nil {
				log.Error("Prove new blocks error", "error", err)
			}
		case e := <-p.blockVerifiedCh:
			if err := p.onBlockVerified(p.ctx, e); err != nil {
				log.Error("Handle BlockVerified event error", "error", err)
			}
		case e := <-p.transitionProvedCh:
			if err := p.onTransitionProved(p.ctx, e); err != nil {
				log.Error("Handle TransitionProved event error", "error", err)
			}
		case e := <-p.proofWindowExpiredCh:
			if err := p.onProvingWindowExpired(p.ctx, e); err != nil {
				log.Error("Handle provingWindow expired event error", "error", err)
			}
		case <-p.blockProposedCh:
			reqProving()
		case <-forceProvingTicker.C:
			reqProving()
		}
	}
}

// Close closes the prover instance.
func (p *Prover) Close(ctx context.Context) {
	p.closeSubscription()
	if err := p.srv.Shutdown(ctx); err != nil {
		log.Error("Failed to shut down prover server", "error", err)
	}
	p.wg.Wait()
}

// proveOp performs a proving operation, find current unproven blocks, then
// request generating proofs for them.
func (p *Prover) proveOp() error {
	firstTry := true

	for firstTry || p.reorgDetectedFlag {
		p.reorgDetectedFlag = false
		firstTry = false

		iter, err := eventIterator.NewBlockProposedIterator(p.ctx, &eventIterator.BlockProposedIteratorConfig{
			Client:               p.rpc.L1,
			TaikoL1:              p.rpc.TaikoL1,
			StartHeight:          new(big.Int).SetUint64(p.l1Current.Number.Uint64()),
			OnBlockProposedEvent: p.onBlockProposed,
		})
		if err != nil {
			log.Error("Failed to start event iterator", "event", "BlockProposed", "error", err)
			return err
		}

		if err := iter.Iter(); err != nil {
			return err
		}
	}

	return nil
}

// onBlockProposed tries to prove that the newly proposed block is valid/invalid.
func (p *Prover) onBlockProposed(
	ctx context.Context,
	event *bindings.TaikoL1ClientBlockProposed,
	end eventIterator.EndBlockProposedEventIterFunc,
) error {
	// If there is newly generated proofs, we need to submit them as soon as possible.
	if len(p.proofGenerationCh) > 0 {
		log.Info("onBlockProposed early return", "proofGenerationChannelLength", len(p.proofGenerationCh))

		end()
		return nil
	}

	if _, err := p.rpc.WaitL1Origin(ctx, event.BlockId); err != nil {
		return fmt.Errorf("failed to wait L1Origin (eventID %d): %w", event.BlockId, err)
	}

	// Check whether the L2 EE's anchored L1 info, to see if the L1 chain has been reorged.
	reorged, l1CurrentToReset, lastHandledBlockIDToReset, err := p.rpc.CheckL1ReorgFromL2EE(
		ctx,
		new(big.Int).Sub(event.BlockId, common.Big1),
	)
	if err != nil {
		return fmt.Errorf("failed to check whether L1 chain was reorged from L2EE (eventID %d): %w", event.BlockId, err)
	}

	// Then check the l1Current cursor at first, to see if the L1 chain has been reorged.
	if !reorged {
		if reorged, l1CurrentToReset, lastHandledBlockIDToReset, err = p.rpc.CheckL1ReorgFromL1Cursor(
			ctx,
			p.l1Current,
			p.genesisHeightL1,
		); err != nil {
			return fmt.Errorf(
				"failed to check whether L1 chain was reorged from l1Current (eventID %d): %w",
				event.BlockId,
				err,
			)
		}
	}

	if reorged {
		log.Info(
			"Reset L1Current cursor due to reorg",
			"l1CurrentHeightOld", p.l1Current,
			"l1CurrentHeightNew", l1CurrentToReset.Number,
			"lastHandledBlockIDOld", p.lastHandledBlockID,
			"lastHandledBlockIDNew", lastHandledBlockIDToReset,
		)
		p.l1Current = l1CurrentToReset
		if lastHandledBlockIDToReset == nil {
			p.lastHandledBlockID = 0
		} else {
			p.lastHandledBlockID = lastHandledBlockIDToReset.Uint64()
		}
		p.reorgDetectedFlag = true
		end()
		return nil
	}

	if event.BlockId.Uint64() <= p.lastHandledBlockID {
		return nil
	}

	currentL1OriginHeader, err := p.rpc.L1.HeaderByNumber(ctx, new(big.Int).SetUint64(event.Meta.L1Height))
	if err != nil {
		return fmt.Errorf("failed to get L1 header, height %d: %w", event.Meta.L1Height, err)
	}

	if currentL1OriginHeader.Hash() != event.Meta.L1Hash {
		log.Warn(
			"L1 block hash mismatch due to L1 reorg",
			"height", event.Meta.L1Height,
			"currentL1OriginHeader", currentL1OriginHeader.Hash(),
			"L1HashInEvent", event.Meta.L1Hash,
		)

		return fmt.Errorf(
			"L1 block hash mismatch due to L1 reorg: %s != %s",
			currentL1OriginHeader.Hash(),
			event.Meta.L1Hash,
		)
	}

	log.Info(
		"Proposed block",
		"L1Height", event.Raw.BlockNumber,
		"L1Hash", event.Raw.BlockHash,
		"BlockID", event.BlockId,
		"Removed", event.Raw.Removed,
	)
	metrics.ProverReceivedProposedBlockGauge.Update(event.BlockId.Int64())

	handleBlockProposedEvent := func() error {
		defer func() { <-p.proposeConcurrencyGuard }()

		// Check whether the block has been verified.
		isVerified, err := p.isBlockVerified(event.BlockId)
		if err != nil {
			return fmt.Errorf("failed to check if the current L2 block is verified: %w", err)
		}

		if isVerified {
			log.Info("📋 Block has been verified", "blockID", event.BlockId)
			return nil
		}

		// Check whether the block's proof is still needed.
		needNewProof, err := rpc.NeedNewProof(
			p.ctx,
			p.rpc,
			event.BlockId,
			p.proverAddress,
		)
		if err != nil {
			return fmt.Errorf("failed to check whether the L2 block needs a new proof: %w", err)
		}

		if !needNewProof {
			return nil
		}

		log.Info(
			"Proposed block information",
			"blockID", event.BlockId,
			"assignedProver", event.AssignedProver,
			"minTier", event.MinTier,
		)

		provingWindow, err := p.getProvingWindow(event)
		if err != nil {
			return fmt.Errorf("failed to get proving window: %w", err)
		}

		var (
			now                    = uint64(time.Now().Unix())
			provingWindowExpiresAt = currentL1OriginHeader.Time + uint64(provingWindow.Seconds())
			provingWindowExpired   = now > provingWindowExpiresAt
			timeToExpire           = time.Duration(provingWindowExpiresAt-now) * time.Second
		)
		if event.AssignedProver != p.proverAddress && !provingWindowExpired {
			log.Info(
				"Proposed block is not provable",
				"blockID", event.BlockId,
				"prover", event.AssignedProver,
				"expiresAt", provingWindowExpiresAt,
				"timeToExpire", timeToExpire,
			)

			if p.cfg.ProveUnassignedBlocks {
				log.Info("Add proposed block to wait for proof window expiration", "blockID", event.BlockId)
				time.AfterFunc(timeToExpire, func() { p.proofWindowExpiredCh <- event })
			}

			return nil
		}

		// If set not to prove unassigned blocks, this block is still not provable
		// by us even though its open proving.
		if provingWindowExpired && event.AssignedProver == p.proverAddress && !p.cfg.ProveUnassignedBlocks {
			log.Info("Skip proving expired blocks", "blockID", event.BlockId)
			return nil
		}

		log.Info(
			"Proposed block is provable",
			"blockID", event.BlockId,
			"prover", event.AssignedProver,
			"expiresAt", provingWindowExpiresAt,
			"minTier", event.MinTier,
		)

		metrics.ProverProofsAssigned.Inc(1)

		if !p.cfg.GuardianProver {
			if _, ok := p.capacityManager.TakeOneCapacity(event.BlockId.Uint64()); !ok {
				return errNoCapacity
			}
		}

		if proofSubmitter := p.selectSubmitter(event.MinTier); proofSubmitter != nil {
			return proofSubmitter.RequestProof(ctx, event)
		}

		return nil
	}

	p.proposeConcurrencyGuard <- struct{}{}

	newL1Current, err := p.rpc.L1.HeaderByHash(ctx, event.Raw.BlockHash)
	if err != nil {
		return err
	}
	p.l1Current = newL1Current
	p.lastHandledBlockID = event.BlockId.Uint64()

	go func() {
		if err := backoff.Retry(
			func() error { return handleBlockProposedEvent() },
			backoff.WithMaxRetries(backoff.NewConstantBackOff(p.cfg.BackOffRetryInterval), p.cfg.BackOffMaxRetrys),
		); err != nil {
			log.Error("Handle new BlockProposed event error", "error", err)
		}
	}()

	return nil
}

// submitProofOp performs a proof submission operation.
func (p *Prover) submitProofOp(ctx context.Context, proofWithHeader *proofProducer.ProofWithHeader) {
	p.submitProofConcurrencyGuard <- struct{}{}

	go func() {
		defer func() {
			<-p.submitProofConcurrencyGuard
			if !p.cfg.GuardianProver {
				_, released := p.capacityManager.ReleaseOneCapacity(proofWithHeader.Meta.Id)
				if !released {
					log.Error("Failed to release capacity", "id", proofWithHeader.Meta.Id)
				}
			}
		}()

		if err := backoff.Retry(
			func() error {
				proofSubmitter := p.getSubmitterByTier(proofWithHeader.Tier)
				if proofSubmitter == nil {
					return nil
				}

				if err := proofSubmitter.SubmitProof(p.ctx, proofWithHeader); err != nil {
					log.Error("Submit proof error", "error", err)
					return err
				}

				return nil
			},
			backoff.WithMaxRetries(backoff.NewConstantBackOff(p.cfg.BackOffRetryInterval), p.cfg.BackOffMaxRetrys),
		); err != nil {
			log.Error("Submit proof error", "error", err)
		}
	}()
}

// onBlockVerified update the latestVerified block in current state, and cancels
// the block being proven if it's verified.
func (p *Prover) onBlockVerified(ctx context.Context, event *bindings.TaikoL1ClientBlockVerified) error {
	metrics.ProverLatestVerifiedIDGauge.Update(event.BlockId.Int64())

	p.latestVerifiedL1Height = event.Raw.BlockNumber

	log.Info(
		"New verified block",
		"blockID", event.BlockId,
		"hash", common.BytesToHash(event.BlockHash[:]),
		"prover", event.Prover,
	)

	return nil
}

// Name returns the application name.
func (p *Prover) Name() string {
	return "prover"
}

// initL1Current initializes prover's L1Current cursor.
func (p *Prover) initL1Current(startingBlockID *big.Int) error {
	if err := p.rpc.WaitTillL2ExecutionEngineSynced(p.ctx); err != nil {
		return err
	}

	stateVars, err := p.rpc.GetProtocolStateVariables(&bind.CallOpts{Context: p.ctx})
	if err != nil {
		return err
	}
	p.genesisHeightL1 = stateVars.GenesisHeight

	if startingBlockID == nil {
		if stateVars.LastVerifiedBlockId == 0 {
			genesisL1Header, err := p.rpc.L1.HeaderByNumber(p.ctx, new(big.Int).SetUint64(stateVars.GenesisHeight))
			if err != nil {
				return err
			}

			p.l1Current = genesisL1Header
			return nil
		}

		startingBlockID = new(big.Int).SetUint64(stateVars.LastVerifiedBlockId)
	}

	log.Info("Init L1Current cursor", "startingBlockID", startingBlockID)

	latestVerifiedHeaderL1Origin, err := p.rpc.L2.L1OriginByID(p.ctx, startingBlockID)
	if err != nil {
		if err.Error() == ethereum.NotFound.Error() {
			log.Warn("Failed to find L1Origin for blockID, use latest L1 head instead", "blockID", startingBlockID)
			l1Head, err := p.rpc.L1.HeaderByNumber(p.ctx, nil)
			if err != nil {
				return err
			}

			p.l1Current = l1Head
			return nil
		}
		return err
	}

	if p.l1Current, err = p.rpc.L1.HeaderByHash(p.ctx, latestVerifiedHeaderL1Origin.L1BlockHash); err != nil {
		return err
	}

	return nil
}

// isBlockVerified checks whether the given block has been verified by other provers.
func (p *Prover) isBlockVerified(id *big.Int) (bool, error) {
	stateVars, err := p.rpc.GetProtocolStateVariables(&bind.CallOpts{Context: p.ctx})
	if err != nil {
		return false, err
	}

	return id.Uint64() <= stateVars.LastVerifiedBlockId, nil
}

// initSubscription initializes all subscriptions in current prover instance.
func (p *Prover) initSubscription() {
	p.blockProposedSub = rpc.SubscribeBlockProposed(p.rpc.TaikoL1, p.blockProposedCh)
	p.blockVerifiedSub = rpc.SubscribeBlockVerified(p.rpc.TaikoL1, p.blockVerifiedCh)
	p.transitionProvedSub = rpc.SubscribeTransitionProved(p.rpc.TaikoL1, p.transitionProvedCh)
}

// closeSubscription closes all subscriptions.
func (p *Prover) closeSubscription() {
	p.blockVerifiedSub.Unsubscribe()
	p.blockProposedSub.Unsubscribe()
}

// isValidProof checks if the given proof is a valid one.
func (p *Prover) isValidProof(
	ctx context.Context,
	event *bindings.TaikoL1ClientTransitionProved,
) (bool, error) {
	parent, err := p.rpc.L2ParentByBlockId(ctx, event.BlockId)
	if err != nil {
		return false, err
	}

	block, err := p.rpc.L2.BlockByNumber(ctx, event.BlockId)
	if err != nil {
		return false, err
	}

	l2SignalService, err := p.rpc.TaikoL2.Resolve0(nil, rpc.StringToBytes32("signal_service"), false)
	if err != nil {
		return false, err
	}

	signalRoot, err := p.rpc.GetStorageRoot(
		ctx,
		p.rpc.L2GethClient,
		l2SignalService,
		event.BlockId,
	)
	if err != nil {
		return false, err
	}

	return parent.Hash() == event.ParentHash &&
		block.Hash() == event.BlockHash &&
		signalRoot == event.SignalRoot, nil
}

// onTransitionProved verifies the proven block hash and will try contesting it if the block hash is wrong.
func (p *Prover) onTransitionProved(ctx context.Context, event *bindings.TaikoL1ClientTransitionProved) error {
	metrics.ProverReceivedProvenBlockGauge.Update(event.BlockId.Int64())

	// If the proof generation is cancellable, cancel it and release the capacity.
	proofSubmitter := p.getSubmitterByTier(event.Tier)
	if proofSubmitter != nil && proofSubmitter.Producer().Cancellable() {
		if err := proofSubmitter.Producer().Cancel(ctx, event.BlockId); err != nil {
			return err
		}
		// No need to check if the release is successful here, since this L2 block might
		// be assigned to other provers.
		p.capacityManager.ReleaseOneCapacity(event.BlockId.Uint64())
	}

	// If this prover is in contest mode, we check the validity of this proof and if it's invalid,
	// contest it with a higher tier proof.
	if !p.cfg.ContestControversialProofs {
		return nil
	}

	isValidProof, err := p.isValidProof(
		ctx,
		event,
	)
	if err != nil {
		return err
	}

	if isValidProof {
		return nil
	}

	l1Height, err := p.rpc.TaikoL2.LatestSyncedL1Height(&bind.CallOpts{Context: ctx, BlockNumber: event.BlockId})
	if err != nil {
		return err
	}

	log.Info(
		"Contest a proven block",
		"blockID", event.BlockId,
		"l1Height", l1Height,
		"tier", event.Tier,
		"parentHash", common.Bytes2Hex(event.ParentHash[:]),
		"blockHash", common.Bytes2Hex(event.BlockHash[:]),
		"signalRoot", common.Bytes2Hex(event.SignalRoot[:]),
	)

	return p.requestProofByBlockID(event.BlockId, new(big.Int).SetUint64(l1Height+1), true)
}

// requestProofByBlockID performs a proving operation for the given block.
func (p *Prover) requestProofByBlockID(blockID *big.Int, l1Height *big.Int, contestOldProof bool) error {
	onBlockProposed := func(
		ctx context.Context,
		event *bindings.TaikoL1ClientBlockProposed,
		end eventIterator.EndBlockProposedEventIterFunc,
	) error {
		// Only filter for exact blockID we want
		if event.BlockId.Cmp(blockID) != 0 {
			return nil
		}

		// Check whether the block has been verified.
		isVerified, err := p.isBlockVerified(event.BlockId)
		if err != nil {
			return fmt.Errorf("failed to check if the current L2 block is verified: %w", err)
		}

		if isVerified {
			log.Info("📋 Block has been verified", "blockID", event.BlockId)
			return nil
		}

		// Make sure to take a capacity before requesting proof
		if !p.cfg.GuardianProver {
			if _, ok := p.capacityManager.TakeOneCapacity(event.BlockId.Uint64()); !ok {
				return errNoCapacity
			}
		}

		p.proposeConcurrencyGuard <- struct{}{}

		var minTier = event.MinTier
		if contestOldProof {
			if p.cfg.GuardianProver {
				minTier = encoding.TierGuardianID
			} else {
				minTier += 1
			}
		}

		if proofSubmitter := p.selectSubmitter(minTier); proofSubmitter != nil {
			return proofSubmitter.RequestProof(ctx, event)
		}

		return nil
	}

	handleBlockProposedEvent := func() error {
		defer func() { <-p.proposeConcurrencyGuard }()
		l1Head, err := p.rpc.L1.BlockNumber(p.ctx)
		if err != nil {
			log.Error("Failed to get L1 block head", "error", err)
			return err
		}
		end := new(big.Int).Add(l1Height, common.Big1)
		if end.Uint64() > l1Head {
			end = new(big.Int).SetUint64(l1Head)
		}
		iter, err := eventIterator.NewBlockProposedIterator(p.ctx, &eventIterator.BlockProposedIteratorConfig{
			Client:               p.rpc.L1,
			TaikoL1:              p.rpc.TaikoL1,
			StartHeight:          new(big.Int).Sub(l1Height, common.Big1),
			EndHeight:            end,
			OnBlockProposedEvent: onBlockProposed,
		})
		if err != nil {
			log.Error("Failed to start event iterator", "event", "BlockProposed", "error", err)
			return err
		}
		return iter.Iter()
	}

	go func() {
		if err := backoff.Retry(
			func() error { return handleBlockProposedEvent() },
			backoff.WithMaxRetries(backoff.NewConstantBackOff(p.cfg.BackOffRetryInterval), p.cfg.BackOffMaxRetrys),
		); err != nil {
			log.Error("Failed to request proof with a given block ID", "blockID", blockID, "error", err)
		}
	}()

	return nil
}

// onProvingWindowExpired tries to submit a proof for an expired block.
func (p *Prover) onProvingWindowExpired(ctx context.Context, e *bindings.TaikoL1ClientBlockProposed) error {
	log.Info("Block proving window is expired", "blockID", e.BlockId, "l1Height", e.Raw.BlockNumber)

	needNewProof, err := rpc.NeedNewProof(ctx, p.rpc, e.BlockId, p.proverAddress)
	if err != nil {
		return err
	}

	if !needNewProof {
		return nil
	}

	return p.requestProofByBlockID(e.BlockId, new(big.Int).SetUint64(e.Raw.BlockNumber), false)
}

// getProvingWindow returns the provingWindow of the given proposed block.
func (p *Prover) getProvingWindow(e *bindings.TaikoL1ClientBlockProposed) (time.Duration, error) {
	for _, t := range p.tiers {
		if e.MinTier == t.ID {
			return time.Duration(t.ProvingWindow) * time.Second, nil
		}
	}

	return 0, errTierNotFound
}

// selectSubmitter returns the proof submitter with the given minTier.
func (p *Prover) selectSubmitter(minTier uint16) proofSubmitter.Submitter {
	for _, s := range p.proofSubmitters {
		if s.Tier() >= minTier {
			return s
		}
	}

	log.Warn("No proof producer / submitter found for the given minTier", "minTier", minTier)

	return nil
}

// getSubmitterByTier returns the proof submitter with the given tier.
func (p *Prover) getSubmitterByTier(tier uint16) proofSubmitter.Submitter {
	for _, s := range p.proofSubmitters {
		if s.Tier() == tier {
			return s
		}
	}

	log.Warn("No proof producer / submitter found for the given tier", "tier", tier)

	return nil
}
