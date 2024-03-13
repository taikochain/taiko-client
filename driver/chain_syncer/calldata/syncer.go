package calldata

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/taikoxyz/taiko-client/bindings"
	"github.com/taikoxyz/taiko-client/bindings/encoding"
	anchorTxConstructor "github.com/taikoxyz/taiko-client/driver/anchor_tx_constructor"
	"github.com/taikoxyz/taiko-client/driver/chain_syncer/beaconsync"
	"github.com/taikoxyz/taiko-client/driver/state"
	txlistfetcher "github.com/taikoxyz/taiko-client/driver/txlist_fetcher"
	"github.com/taikoxyz/taiko-client/internal/metrics"
	eventIterator "github.com/taikoxyz/taiko-client/pkg/chain_iterator/event_iterator"
	"github.com/taikoxyz/taiko-client/pkg/rpc"
	txListValidator "github.com/taikoxyz/taiko-client/pkg/txlist_validator"
)

// Syncer responsible for letting the L2 execution engine catching up with protocol's latest
// pending block through deriving L1 calldata.
type Syncer struct {
	ctx               context.Context
	rpc               *rpc.Client
	state             *state.State
	progressTracker   *beaconsync.SyncProgressTracker          // Sync progress tracker
	anchorConstructor *anchorTxConstructor.AnchorTxConstructor // TaikoL2.anchor transactions constructor
	txListValidator   *txListValidator.TxListValidator         // Transactions list validator
	// Used by BlockInserter
	lastInsertedBlockID *big.Int
	reorgDetectedFlag   bool
}

// NewSyncer creates a new syncer instance.
func NewSyncer(
	ctx context.Context,
	rpc *rpc.Client,
	state *state.State,
	progressTracker *beaconsync.SyncProgressTracker,
) (*Syncer, error) {
	configs, err := rpc.TaikoL1.GetConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("failed to get protocol configs: %w", err)
	}

	constructor, err := anchorTxConstructor.New(rpc)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize anchor constructor: %w", err)
	}

	return &Syncer{
		ctx:               ctx,
		rpc:               rpc,
		state:             state,
		progressTracker:   progressTracker,
		anchorConstructor: constructor,
		txListValidator: txListValidator.NewTxListValidator(
			uint64(configs.BlockMaxGasLimit),
			configs.BlockMaxTxListBytes.Uint64(),
			rpc.L2.ChainID,
		),
	}, nil
}

// ProcessL1Blocks fetches all `TaikoL1.BlockProposed` events between given
// L1 block heights, and then tries inserting them into L2 execution engine's blockchain.
func (s *Syncer) ProcessL1Blocks(ctx context.Context, l1End *types.Header) error {
	firstTry := true
	for firstTry || s.reorgDetectedFlag {
		s.reorgDetectedFlag = false
		firstTry = false

		startL1Current := s.state.GetL1Current()
		// If there is a L1 reorg, sometimes this will happen.
		if startL1Current.Number.Uint64() >= l1End.Number.Uint64() && startL1Current.Hash() != l1End.Hash() {
			newL1Current, err := s.rpc.L1.HeaderByNumber(ctx, new(big.Int).Sub(l1End.Number, common.Big1))
			if err != nil {
				return err
			}

			log.Info(
				"Reorg detected",
				"oldL1CurrentHeight", startL1Current.Number,
				"oldL1CurrentHash", startL1Current.Hash(),
				"newL1CurrentHeight", newL1Current.Number,
				"newL1CurrentHash", newL1Current.Hash(),
				"l1Head", l1End.Number,
			)

			s.state.SetL1Current(newL1Current)
			s.lastInsertedBlockID = nil
		}

		iter, err := eventIterator.NewBlockProposedIterator(ctx, &eventIterator.BlockProposedIteratorConfig{
			Client:               s.rpc.L1,
			TaikoL1:              s.rpc.TaikoL1,
			StartHeight:          s.state.GetL1Current().Number,
			EndHeight:            l1End.Number,
			FilterQuery:          nil,
			OnBlockProposedEvent: s.onBlockProposed,
		})
		if err != nil {
			return err
		}

		if err := iter.Iter(); err != nil {
			return err
		}
	}

	s.state.SetL1Current(l1End)
	metrics.DriverL1CurrentHeightGauge.Update(s.state.GetL1Current().Number.Int64())

	return nil
}

// OnBlockProposed is a `BlockProposed` event callback which responsible for
// inserting the proposed block one by one to the L2 execution engine.
func (s *Syncer) onBlockProposed(
	ctx context.Context,
	event *bindings.TaikoL1ClientBlockProposed,
	endIter eventIterator.EndBlockProposedEventIterFunc,
) error {
	// We simply ignore the genesis block's `BlockProposed` event.
	if event.BlockId.Cmp(common.Big0) == 0 {
		return nil
	}

	// If we are not inserting a block whose parent block is the latest verified block in protocol,
	// and the node hasn't just finished the P2P sync, we check if the L1 chain has been reorged.
	if !s.progressTracker.Triggered() {
		reorgCheckResult, err := s.checkReorg(ctx, event)
		if err != nil {
			return err
		}

		if reorgCheckResult.IsReorged {
			log.Info(
				"Reset L1Current cursor due to L1 reorg",
				"l1CurrentHeightOld", s.state.GetL1Current().Number,
				"l1CurrentHashOld", s.state.GetL1Current().Hash(),
				"l1CurrentHeightNew", reorgCheckResult.L1CurrentToReset.Number,
				"l1CurrentHashNew", reorgCheckResult.L1CurrentToReset.Hash(),
				"lastInsertedBlockIDOld", s.lastInsertedBlockID,
				"lastInsertedBlockIDNew", reorgCheckResult.LastHandledBlockIDToReset,
			)
			s.state.SetL1Current(reorgCheckResult.L1CurrentToReset)
			s.lastInsertedBlockID = reorgCheckResult.LastHandledBlockIDToReset
			s.reorgDetectedFlag = true
			endIter()

			return nil
		}
	}
	// Ignore those already inserted blocks.
	if s.lastInsertedBlockID != nil && event.BlockId.Cmp(s.lastInsertedBlockID) <= 0 {
		return nil
	}

	log.Info(
		"New BlockProposed event",
		"l1Height", event.Raw.BlockNumber,
		"l1Hash", event.Raw.BlockHash,
		"blockID", event.BlockId,
		"removed", event.Raw.Removed,
	)

	// Fetch the L2 parent block.
	var (
		parent *types.Header
		err    error
	)
	if s.progressTracker.Triggered() {
		// Already synced through beacon sync, just skip this event.
		if event.BlockId.Cmp(s.progressTracker.LastSyncedVerifiedBlockID()) <= 0 {
			return nil
		}

		parent, err = s.rpc.L2.HeaderByHash(ctx, s.progressTracker.LastSyncedVerifiedBlockHash())
	} else {
		parent, err = s.rpc.L2ParentByBlockID(ctx, event.BlockId)
	}

	if err != nil {
		return fmt.Errorf("failed to fetch L2 parent block: %w", err)
	}

	log.Debug("Parent block", "height", parent.Number, "hash", parent.Hash())

	tx, err := s.rpc.L1.TransactionInBlock(ctx, event.Raw.BlockHash, event.Raw.TxIndex)
	if err != nil {
		return fmt.Errorf("failed to fetch original TaikoL1.proposeBlock transaction: %w", err)
	}

	var txListDecoder txlistfetcher.TxListFetcher
	if event.Meta.BlobUsed {
		txListDecoder = txlistfetcher.NewBlobTxListFetcher(s.rpc)
	} else {
		txListDecoder = new(txlistfetcher.CalldataFetcher)
	}
	txListBytes, err := txListDecoder.Fetch(ctx, tx, &event.Meta)
	if err != nil {
		if errors.Is(err, rpc.ErrBlobInvalid) {
			log.Info("Invalid blob detected", "blockID", event.BlockId)
			txListBytes = []byte{}
		} else {
			return fmt.Errorf("failed to decode tx list: %w", err)
		}
	}

	l1Origin := &rawdb.L1Origin{
		BlockID:       event.BlockId,
		L2BlockHash:   common.Hash{}, // Will be set by taiko-geth.
		L1BlockHeight: new(big.Int).SetUint64(event.Raw.BlockNumber),
		L1BlockHash:   event.Raw.BlockHash,
	}

	if event.Meta.Timestamp > uint64(time.Now().Unix()) {
		log.Warn("Future L2 block, waiting", "L2BlockTimestamp", event.Meta.Timestamp, "now", time.Now().Unix())
		time.Sleep(time.Until(time.Unix(int64(event.Meta.Timestamp), 0)))
	}

	// If the transactions list is invalid, we simply insert an empty L2 block.
	if !s.txListValidator.ValidateTxList(event.BlockId, txListBytes, event.Meta.BlobUsed) {
		log.Info("Invalid transactions list, insert an empty L2 block instead", "blockID", event.BlockId)
		txListBytes = []byte{}
	}

	payloadData, err := s.insertNewHead(
		ctx,
		event,
		parent,
		s.state.GetHeadBlockID(),
		txListBytes,
		l1Origin,
	)
	if err != nil {
		return fmt.Errorf("failed to insert new head to L2 execution engine: %w", err)
	}

	log.Debug("Payload data", "hash", payloadData.BlockHash, "txs", len(payloadData.Transactions))

	log.Info(
		"🔗 New L2 block inserted",
		"blockID", event.BlockId,
		"height", payloadData.Number,
		"hash", payloadData.BlockHash,
		"transactions", len(payloadData.Transactions),
		"baseFee", payloadData.BaseFeePerGas,
		"withdrawals", len(payloadData.Withdrawals),
	)

	metrics.DriverL1CurrentHeightGauge.Update(int64(event.Raw.BlockNumber))
	s.lastInsertedBlockID = event.BlockId

	if s.progressTracker.Triggered() {
		s.progressTracker.ClearMeta()
	}

	return nil
}

// insertNewHead tries to insert a new head block to the L2 execution engine's local
// block chain through Engine APIs.
func (s *Syncer) insertNewHead(
	ctx context.Context,
	event *bindings.TaikoL1ClientBlockProposed,
	parent *types.Header,
	headBlockID *big.Int,
	txListBytes []byte,
	l1Origin *rawdb.L1Origin,
) (*engine.ExecutableData, error) {
	log.Debug(
		"Try to insert a new L2 head block",
		"parentNumber", parent.Number,
		"parentHash", parent.Hash(),
		"headBlockID", headBlockID,
		"l1Origin", l1Origin,
	)

	// Insert a TaikoL2.anchor transaction at transactions list head
	var txList []*types.Transaction
	if len(txListBytes) != 0 {
		if err := rlp.DecodeBytes(txListBytes, &txList); err != nil {
			log.Error("Invalid txList bytes", "blockID", event.BlockId)
			return nil, err
		}
	}

	// Get L2 baseFee
	baseFee, err := s.rpc.TaikoL2.GetBasefee(
		&bind.CallOpts{BlockNumber: parent.Number, Context: ctx},
		event.Meta.L1Height,
		uint32(parent.GasUsed),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get L2 baseFee: %w", encoding.TryParsingCustomError(err))
	}

	log.Info(
		"L2 baseFee",
		"blockID", event.BlockId,
		"baseFee", baseFee,
		"syncedL1Height", event.Meta.L1Height,
		"parentGasUsed", parent.GasUsed,
	)

	// Get withdrawals
	withdrawals := make(types.Withdrawals, len(event.DepositsProcessed))
	for i, d := range event.DepositsProcessed {
		withdrawals[i] = &types.Withdrawal{Address: d.Recipient, Amount: d.Amount.Uint64(), Index: d.Id}
	}

	// Assemble a TaikoL2.anchor transaction
	anchorTx, err := s.anchorConstructor.AssembleAnchorTx(
		ctx,
		new(big.Int).SetUint64(event.Meta.L1Height),
		event.Meta.L1Hash,
		new(big.Int).Add(parent.Number, common.Big1),
		baseFee,
		parent.GasUsed,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create TaikoL2.anchor transaction: %w", err)
	}

	txList = append([]*types.Transaction{anchorTx}, txList...)
	if txListBytes, err = rlp.EncodeToBytes(txList); err != nil {
		log.Error("Encode txList error", "blockID", event.BlockId, "error", err)
		return nil, err
	}

	payload, err := s.createExecutionPayloads(
		ctx,
		event,
		parent.Hash(),
		l1Origin,
		headBlockID,
		txListBytes,
		baseFee,
		withdrawals,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution payloads: %w", err)
	}

	fc := &engine.ForkchoiceStateV1{HeadBlockHash: parent.Hash()}

	// Update the fork choice
	fc.HeadBlockHash = payload.BlockHash
	fcRes, err := s.rpc.L2Engine.ForkchoiceUpdate(ctx, fc, nil)
	if err != nil {
		return nil, err
	}
	if fcRes.PayloadStatus.Status != engine.VALID {
		return nil, fmt.Errorf("unexpected ForkchoiceUpdate response status: %s", fcRes.PayloadStatus.Status)
	}

	return payload, nil
}

// createExecutionPayloads creates a new execution payloads through
// Engine APIs.
func (s *Syncer) createExecutionPayloads(
	ctx context.Context,
	event *bindings.TaikoL1ClientBlockProposed,
	parentHash common.Hash,
	l1Origin *rawdb.L1Origin,
	headBlockID *big.Int,
	txListBytes []byte,
	baseFee *big.Int,
	withdrawals types.Withdrawals,
) (payloadData *engine.ExecutableData, err error) {
	fc := &engine.ForkchoiceStateV1{HeadBlockHash: parentHash}
	attributes := &engine.PayloadAttributes{
		Timestamp:             event.Meta.Timestamp,
		Random:                event.Meta.Difficulty,
		SuggestedFeeRecipient: event.Meta.Coinbase,
		Withdrawals:           withdrawals,
		BlockMetadata: &engine.BlockMetadata{
			HighestBlockID: headBlockID,
			Beneficiary:    event.Meta.Coinbase,
			GasLimit:       uint64(event.Meta.GasLimit) + anchorTxConstructor.AnchorGasLimit,
			Timestamp:      event.Meta.Timestamp,
			TxList:         txListBytes,
			MixHash:        event.Meta.Difficulty,
			ExtraData:      event.Meta.ExtraData[:],
		},
		BaseFeePerGas: baseFee,
		L1Origin:      l1Origin,
	}

	log.Debug(
		"PayloadAttributes",
		"blockID", event.BlockId,
		"timestamp", attributes.Timestamp,
		"random", attributes.Random,
		"suggestedFeeRecipient", attributes.SuggestedFeeRecipient,
		"withdrawals", len(attributes.Withdrawals),
		"highestBlockID", attributes.BlockMetadata.HighestBlockID,
		"gasLimit", attributes.BlockMetadata.GasLimit,
		"timestamp", attributes.BlockMetadata.Timestamp,
		"mixHash", attributes.BlockMetadata.MixHash,
		"baseFee", attributes.BaseFeePerGas,
		"extraData", string(attributes.BlockMetadata.ExtraData),
		"l1OriginHeight", attributes.L1Origin.L1BlockHeight,
		"l1OriginHash", attributes.L1Origin.L1BlockHash,
	)

	// Step 1, prepare a payload
	fcRes, err := s.rpc.L2Engine.ForkchoiceUpdate(ctx, fc, attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to update fork choice: %w", err)
	}
	if fcRes.PayloadStatus.Status != engine.VALID {
		return nil, fmt.Errorf("unexpected ForkchoiceUpdate response status: %s", fcRes.PayloadStatus.Status)
	}
	if fcRes.PayloadID == nil {
		return nil, errors.New("empty payload ID")
	}

	// Step 2, get the payload
	payload, err := s.rpc.L2Engine.GetPayload(ctx, fcRes.PayloadID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payload: %w", err)
	}

	log.Debug(
		"Payload",
		"blockID", event.BlockId,
		"baseFee", payload.BaseFeePerGas,
		"number", payload.Number,
		"hash", payload.BlockHash,
		"gasLimit", payload.GasLimit,
		"gasUsed", payload.GasUsed,
		"timestamp", payload.Timestamp,
		"withdrawalsHash", payload.WithdrawalsHash,
	)

	// Step 3, execute the payload
	execStatus, err := s.rpc.L2Engine.NewPayload(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new payload: %w", err)
	}
	if execStatus.Status != engine.VALID {
		return nil, fmt.Errorf("unexpected NewPayload response status: %s", execStatus.Status)
	}

	return payload, nil
}

// checkLastVerifiedBlockMismatch checks if there is a mismatch between protocol's last verified block hash and
// the corresponding L2 EE block hash.
func (s *Syncer) checkLastVerifiedBlockMismatch(ctx context.Context) (bool, error) {
	stateVars, err := s.rpc.GetProtocolStateVariables(&bind.CallOpts{Context: ctx})
	if err != nil {
		return false, err
	}

	if s.state.GetL2Head().Number.Uint64() < stateVars.B.LastVerifiedBlockId {
		return false, nil
	}

	blockInfo, err := s.rpc.TaikoL1.GetBlock(&bind.CallOpts{Context: ctx}, stateVars.B.LastVerifiedBlockId)
	if err != nil {
		return false, err
	}

	l2Header, err := s.rpc.L2.HeaderByNumber(ctx, new(big.Int).SetUint64(stateVars.B.LastVerifiedBlockId))
	if err != nil {
		return false, err
	}

	return blockInfo.Ts.BlockHash != l2Header.Hash(), nil
}

// checkReorg checks whether the L1 chain has been reorged, and resets the L1Current cursor if necessary.
func (s *Syncer) checkReorg(
	ctx context.Context,
	event *bindings.TaikoL1ClientBlockProposed,
) (*rpc.ReorgCheckResult, error) {
	var (
		reorgCheckResult = new(rpc.ReorgCheckResult)
		err              error
	)

	// If the L2 chain is at genesis, we don't need to check L1 reorg.
	if s.state.GetL1Current().Number == s.state.GenesisL1Height {
		return reorgCheckResult, nil
	}

	// 1. The latest verified block
	mismatch, err := s.checkLastVerifiedBlockMismatch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if last verified block in L2 EE has been reorged: %w", err)
	}

	// If the latest verified block in chain is mismatched, we reset the L2 chain to genesis, and restart
	// the calldata sync process.
	// TODO(Gavin): improve this approach.
	if mismatch {
		log.Warn("The latest verified block mismatch detected, reset L2 chain to genesis")

		genesisL1Header, err := s.rpc.GetGenesisL1Header(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch genesis L1 header: %w", err)
		}

		reorgCheckResult.IsReorged = true
		reorgCheckResult.L1CurrentToReset = genesisL1Header
		reorgCheckResult.LastHandledBlockIDToReset = common.Big0
	} else {
		// 2. Parent block
		reorgCheckResult, err = s.rpc.CheckL1Reorg(
			ctx,
			new(big.Int).Sub(event.BlockId, common.Big1),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to check whether L1 chain has been reorged: %w", err)
		}
	}

	return reorgCheckResult, nil
}
