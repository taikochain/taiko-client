package driver

import (
	"context"
	"math/big"
	"sync/atomic"

	"github.com/cenkalti/backoff/v4"
	"github.com/taikochain/client-mono/bindings"
	ethereum "github.com/taikochain/taiko-client"
	"github.com/taikochain/taiko-client/common"
	"github.com/taikochain/taiko-client/core/types"
	"github.com/taikochain/taiko-client/event"
	"github.com/taikochain/taiko-client/log"
	"golang.org/x/sync/errgroup"
)

type State struct {
	// Subscriptions, will automatically resubscribe on errors
	l1HeadSub           event.Subscription // L1 new heads
	l2BlockFinalizedSub event.Subscription // TaikoL1.BlockFinalized events
	l2BlockProposedSub  event.Subscription // TaikoL1.BlockProposed events

	// Feeds
	l1HeadsFeed event.Feed // L1 new heads notification feed

	l1Head              *atomic.Value // Latest known L1 head
	l2HeadBlockID       *atomic.Value // Latest known L2 block ID
	l2FinalizedHeadHash *atomic.Value // Latest known L2 finalized head
	L1Current           *types.Header // Current L1 block sync cursor

	// Constants
	anchorTxGasLimit  *big.Int
	maxTxlistBytes    *big.Int
	maxBlockNumTxs    *big.Int
	maxBlocksGasLimit *big.Int
	minTxGasLimit     *big.Int

	rpc *RPCClient // L1/L2 RPC clients
}

// NewState creates a new driver state instance.
func NewState(ctx context.Context, rpc *RPCClient) (*State, error) {
	// Set the L2 head's latest known L1 origin as current L1 sync cursor.
	latestL2KnownL1Header, err := rpc.LatestL2KnownL1Header(ctx)
	if err != nil {
		return nil, err
	}

	_, _, _, _, _,
		maxBlocksGasLimit, maxBlockNumTxs, _, maxTxlistBytes, minTxGasLimit,
		anchorTxGasLimit, _, _, err := rpc.taikoL1.GetConstants(nil)
	if err != nil {
		return nil, err
	}

	log.Info(
		"Taiko constants",
		"anchorTxGasLimit", anchorTxGasLimit,
		"maxTxlistBytes", maxTxlistBytes,
		"maxBlockNumTxs", maxBlockNumTxs,
		"maxBlocksGasLimit", maxBlocksGasLimit,
		"minTxGasLimit", minTxGasLimit,
	)

	s := &State{
		rpc:                 rpc,
		anchorTxGasLimit:    anchorTxGasLimit,
		maxTxlistBytes:      maxTxlistBytes,
		maxBlockNumTxs:      maxBlockNumTxs,
		maxBlocksGasLimit:   maxBlocksGasLimit,
		minTxGasLimit:       minTxGasLimit,
		l1Head:              new(atomic.Value),
		l2HeadBlockID:       new(atomic.Value),
		l2FinalizedHeadHash: new(atomic.Value),
		L1Current:           latestL2KnownL1Header,
	}

	if err := s.initSubscriptions(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

// Close closes all inner subscriptions.
func (s *State) Close() error {
	s.l1HeadSub.Unsubscribe()
	s.l2BlockFinalizedSub.Unsubscribe()
	s.l2BlockProposedSub.Unsubscribe()
	return nil
}

// ConfirmL1Current ensures that the local L1 sync cursor has not been reorged.
func (s *State) ConfirmL1Current(ctx context.Context) (*types.Header, error) {
	// Check whether the L1 sync cursor has been reorged
	l1Current, err := s.rpc.l1.HeaderByHash(ctx, s.L1Current.Hash())
	if err != nil {
		return nil, err
	}

	// Not reorged
	if l1Current != nil {
		return s.L1Current, nil
	}

	// Reorg detected, update sync cursor's height to last saved height minus
	// `reorgRollbackDepth`
	var newL1SyncCursorHeight *big.Int
	if s.L1Current.Number.Uint64() > ReorgRollbackDepth {
		newL1SyncCursorHeight = new(big.Int).SetUint64(
			s.L1Current.Number.Uint64() - ReorgRollbackDepth,
		)
	} else {
		newL1SyncCursorHeight = common.Big0
	}

	s.L1Current, err = s.rpc.l1.HeaderByNumber(ctx, newL1SyncCursorHeight)
	if err != nil {
		return nil, err
	}

	return s.L1Current, nil
}

// initSubscriptions initializes all subscriptions in the given state instance.
func (s *State) initSubscriptions(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error { return s.subL1Head(ctx) })
	g.Go(func() error { return s.subBlockFinalized(ctx) })
	g.Go(func() error { return s.subLastPendingBlockID(ctx) })

	return g.Wait()
}

// subL1Head fetches the latest L1 Head from L1 node, then saves it to current
// driver state, and then start subscribing.
func (s *State) subL1Head(ctx context.Context) error {
	l1Head, err := s.rpc.l1.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}

	s.setL1Head(l1Head)

	s.l1HeadSub = event.ResubscribeErr(
		backoff.DefaultMaxInterval,
		func(ctx context.Context, err error) (event.Subscription, error) {
			if err != nil {
				log.Warn("Failed to subscribe L1 head, try resubscribing", "error", err)
			}

			return s.watchL1Head(ctx)
		},
	)
	go func() {
		err, ok := <-s.l1HeadSub.Err()
		if !ok {
			return
		}
		log.Error("Subscribe L1 head error", "error", err)
	}()

	return nil
}

// watchL1Head watches new L1 head blocks and keep updating current
// driver state.
func (s *State) watchL1Head(ctx context.Context) (event.Subscription, error) {
	newL1HeadCh := make(chan *types.Header, 10)
	sub, err := s.rpc.l1.SubscribeNewHead(context.Background(), newL1HeadCh)
	if err != nil {
		log.Error("Create L1 heads subscription error", "error", err)
		return nil, err
	}

	return event.NewSubscription(func(quitCh <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case newHead := <-newL1HeadCh:
				s.setL1Head(newHead)
				s.l1HeadsFeed.Send(newHead)
			case err := <-sub.Err():
				return err
			case <-quitCh:
				return nil
			}
		}
	}), nil
}

// setL1Head sets the L1 head concurrent safely.
func (s *State) setL1Head(l1Head *types.Header) {
	if l1Head == nil {
		log.Warn("Empty new L1 head")
		return
	}

	log.Info("New L1 head", "height", l1Head.Number, "hash", l1Head.Hash(), "timestamp", l1Head.Time)

	s.l1Head.Store(l1Head)
}

// getL1Head reads the L1 head concurrent safely.
func (s *State) getL1Head() *types.Header {
	return s.l1Head.Load().(*types.Header)
}

// subBlockFinalized fetches the last finalized L2 block hash from
// TaikoL1 contract, saves it to the current driver state, and then start
// subscribing.
func (s *State) subBlockFinalized(ctx context.Context) error {
	_, lastFinalizedHeight, _, _, err := s.rpc.taikoL1.GetStateVariables(nil)
	if err != nil {
		return err
	}

	lastFinalizedBlockHash, err := s.rpc.taikoL1.GetSyncedHeader(
		nil,
		new(big.Int).SetUint64(lastFinalizedHeight),
	)
	if err != nil {
		return err
	}

	s.setLastFinalizedBlockHash(lastFinalizedBlockHash)

	s.l2BlockFinalizedSub = event.ResubscribeErr(
		backoff.DefaultMaxInterval,
		func(ctx context.Context, err error) (event.Subscription, error) {
			if err != nil {
				log.Warn("Failed to subscribe TaikoL1.BlockFinalized events, try resubscribing", "error", err)
			}

			return s.watchBlockFinalized(ctx)
		},
	)
	go func() {
		err, ok := <-s.l2BlockFinalizedSub.Err()
		if !ok {
			return
		}
		log.Error("Subscribe TaikoL1.BlockFinalized error", "error", err)
	}()

	return nil
}

// watchBlockFinalized watches newly finalized blocks and keep updating current
// driver state.
func (s *State) watchBlockFinalized(ctx context.Context) (ethereum.Subscription, error) {
	newBlockFinalizedCh := make(chan *bindings.TaikoL1ClientBlockFinalized, 10)
	sub, err := s.rpc.taikoL1.WatchBlockFinalized(nil, newBlockFinalizedCh, nil)
	if err != nil {
		log.Error("Create TaikoL1.BlockFinalized subscription error", "error", err)
		return nil, err
	}

	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case e := <-newBlockFinalizedCh:
				if e.BlockHash == (common.Hash{}) {
					log.Info("Ignore BlockFinalized event of invalid transaction list", "blockID", e.Id)
					continue
				}
				s.setLastFinalizedBlockHash(e.BlockHash)
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// setLastFinalizedBlockHash sets the last finalized L2 block hash concurrent safely.
func (s *State) setLastFinalizedBlockHash(hash common.Hash) {
	log.Info("New finalized block", "hash", hash)

	s.l2FinalizedHeadHash.Store(hash)
}

// getLastFinalizedBlockHash reads the last finalized L2 block concurrent safely.
func (s *State) getLastFinalizedBlockHash() common.Hash {
	return s.l2FinalizedHeadHash.Load().(common.Hash)
}

// subLastPendingBlockID fetches the last pending block ID from
// taikoL1 contract, saves it to the current driver state, and then starts
// subscribing.
func (s *State) subLastPendingBlockID(ctx context.Context) error {
	_, _, _, nextPendingID, err := s.rpc.taikoL1.GetStateVariables(nil)
	if err != nil {
		return err
	}

	s.setHeadBlockID(new(big.Int).SetUint64(nextPendingID - 1))

	s.l2BlockProposedSub = event.ResubscribeErr(
		backoff.DefaultMaxInterval,
		func(ctx context.Context, err error) (event.Subscription, error) {
			if err != nil {
				log.Warn("Failed to subscribe TaikoL1.BlockProposed events, try resubscribing", "error", err)
			}

			return s.watchBlockProposed(ctx)
		},
	)
	go func() {
		err, ok := <-s.l2BlockProposedSub.Err()
		if !ok {
			return
		}
		log.Error("Subscribe TaikoL1.BlockProposed error", "error", err)
	}()

	return nil
}

// watchBlockProposed watches newly proposed blocks and keeps updating current
// driver state.
func (s *State) watchBlockProposed(ctx context.Context) (ethereum.Subscription, error) {
	newBlockProposedCh := make(chan *bindings.TaikoL1ClientBlockProposed, 10)
	sub, err := s.rpc.taikoL1.WatchBlockProposed(nil, newBlockProposedCh, nil)
	if err != nil {
		log.Error("Create TaikoL1.BlockProposed subscription error", "error", err)
		return nil, err
	}

	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case e := <-newBlockProposedCh:
				s.setHeadBlockID(e.Id)
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// setHeadBlockID sets the last pending block ID concurrent safely.
func (s *State) setHeadBlockID(id *big.Int) {
	log.Info("New head block ID", "ID", id)

	s.l2HeadBlockID.Store(id)
}

// getHeadBlockID reads the last pending block ID concurrent safely.
func (s *State) getHeadBlockID() *big.Int {
	return s.l2HeadBlockID.Load().(*big.Int)
}

// SubL1HeadsFeed registers a subscription of new L1 heads.
func (s *State) SubL1HeadsFeed(ch chan *types.Header) event.Subscription {
	return s.l1HeadsFeed.Subscribe(ch)
}
