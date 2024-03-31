package chainsyncer

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-client/driver/chain_syncer/beaconsync"
	"github.com/taikoxyz/taiko-client/driver/chain_syncer/calldata"
	"github.com/taikoxyz/taiko-client/driver/state"
	"github.com/taikoxyz/taiko-client/pkg/rpc"
)

// L2ChainSyncer is responsible for keeping the L2 execution engine's local chain in sync with the one
// in TaikoL1 contract.
type L2ChainSyncer struct {
	ctx   context.Context
	state *state.State // Driver's state
	rpc   *rpc.Client  // L1/L2 RPC clients

	// Syncers
	beaconSyncer   *beaconsync.Syncer
	calldataSyncer *calldata.Syncer

	// Monitors
	progressTracker *beaconsync.SyncProgressTracker

	// If this flag is activated, will try P2P beacon sync if current node is behind of the protocol's
	// the latest verified block head
	p2pSyncVerifiedBlocks bool
}

// New creates a new chain syncer instance.
func New(
	ctx context.Context,
	rpc *rpc.Client,
	state *state.State,
	p2pSyncVerifiedBlocks bool,
	p2pSyncTimeout time.Duration,
	maxRetrieveExponent uint64,
) (*L2ChainSyncer, error) {
	tracker := beaconsync.NewSyncProgressTracker(rpc.L2, p2pSyncTimeout)
	go tracker.Track(ctx)

	beaconSyncer := beaconsync.NewSyncer(ctx, rpc, state, tracker)
	calldataSyncer, err := calldata.NewSyncer(ctx, rpc, state, tracker, maxRetrieveExponent)
	if err != nil {
		return nil, err
	}

	return &L2ChainSyncer{
		ctx:                   ctx,
		rpc:                   rpc,
		state:                 state,
		beaconSyncer:          beaconSyncer,
		calldataSyncer:        calldataSyncer,
		progressTracker:       tracker,
		p2pSyncVerifiedBlocks: p2pSyncVerifiedBlocks,
	}, nil
}

// Sync performs a sync operation to L2 execution engine's local chain.
func (s *L2ChainSyncer) Sync() error {
	blockID, needNewBeaconSyncTriggered, err := s.needNewBeaconSyncTriggered()
	if err != nil {
		return err
	}
	// If current L2 execution engine's chain is behind of the protocol's latest verified block head, and the
	// `P2PSyncVerifiedBlocks` flag is set, try triggering a beacon sync in L2 execution engine to catch up the
	// latest verified block head.
	if needNewBeaconSyncTriggered {
		if err := s.beaconSyncer.TriggerBeaconSync(blockID); err != nil {
			return fmt.Errorf("trigger beacon sync error: %w", err)
		}

		return nil
	}

	// We have triggered at least a beacon sync in L2 execution engine, we should reset the L1Current
	// cursor at first, before start inserting pending L2 blocks one by one.
	if s.progressTracker.Triggered() {
		log.Info(
			"Switch to insert pending blocks one by one",
			"p2pEnabled", s.p2pSyncVerifiedBlocks,
			"p2pOutOfSync", s.progressTracker.OutOfSync(),
		)

		// Get the execution engine's chain head.
		l2Head, err := s.rpc.L2.HeaderByNumber(s.ctx, nil)
		if err != nil {
			return err
		}

		log.Info(
			"L2 head information",
			"number", l2Head.Number,
			"hash", l2Head.Hash(),
			"lastSyncedVerifiedBlockID", s.progressTracker.LastSyncedVerifiedBlockID(),
			"lastSyncedVerifiedBlockHash", s.progressTracker.LastSyncedVerifiedBlockHash(),
		)

		// Reset the L1Current cursor.
		if err := s.state.ResetL1Current(s.ctx, l2Head.Number); err != nil {
			return err
		}

		// Reset to the latest L2 execution engine's chain status.
		s.progressTracker.UpdateMeta(l2Head.Number.Uint64(), l2Head.Hash())
	}

	// Insert the proposed block one by one.
	return s.calldataSyncer.ProcessL1Blocks(s.ctx)
}

// AheadOfProtocolVerifiedHead checks whether the L2 chain is ahead of verified head in protocol.
func (s *L2ChainSyncer) AheadOfProtocolVerifiedHead(verifiedHeightToCompare uint64) bool {
	log.Debug(
		"Checking whether the execution engine is ahead of protocol's verified head",
		"latestVerifiedBlock", verifiedHeightToCompare,
		"executionEngineHead", s.state.GetL2Head().Number,
	)
	if verifiedHeightToCompare > 0 {
		// If latest verified head height is equal to L2 execution engine's synced head height minus one,
		// we also mark the triggered P2P sync progress as finished to prevent a potential `InsertBlockWithoutSetHead` in
		// execution engine, which may cause errors since we do not pass all transactions in ExecutePayload when calling
		// `NewPayloadV1`.
		verifiedHeightToCompare--
	}

	// If the L2 execution engine's chain is behind of the protocol's latest verified block head,
	// we should keep the beacon sync.
	if s.state.GetL2Head().Number.Uint64() < verifiedHeightToCompare {
		return false
	}

	if s.progressTracker.LastSyncedVerifiedBlockID() != 0 {
		return s.state.GetL2Head().Number.Uint64() >= s.progressTracker.LastSyncedVerifiedBlockID()
	}

	return true
}

// needNewBeaconSyncTriggered checks whether the current L2 execution engine needs to trigger
// another new beacon sync, the following conditions should be met:
// 1. The `P2PSyncVerifiedBlocks` flag is set.
// 2. The protocol's latest verified block head is not zero.
// 3. The L2 execution engine's chain is behind of the protocol's latest verified block head.
// 4. The L2 execution engine's chain have met a sync timeout issue.
func (s *L2ChainSyncer) needNewBeaconSyncTriggered() (uint64, bool, error) {
	// If the flag is not set, we simply return false.
	if !s.p2pSyncVerifiedBlocks {
		return 0, false, nil
	}

	stateVars, err := s.rpc.GetProtocolStateVariables(&bind.CallOpts{Context: s.ctx})
	if err != nil {
		return 0, false, err
	}

	// If the protocol's latest verified block head is zero, we simply return false.
	if stateVars.B.LastVerifiedBlockId == 0 {
		return 0, false, nil
	}

	return stateVars.B.LastVerifiedBlockId, !s.AheadOfProtocolVerifiedHead(stateVars.B.LastVerifiedBlockId) &&
		!s.progressTracker.OutOfSync(), nil
}

// BeaconSyncer returns the inner beacon syncer.
func (s *L2ChainSyncer) BeaconSyncer() *beaconsync.Syncer {
	return s.beaconSyncer
}

// CalldataSyncer returns the inner calldata syncer.
func (s *L2ChainSyncer) CalldataSyncer() *calldata.Syncer {
	return s.calldataSyncer
}
