package driver

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/taikoxyz/taiko-client/bindings"
	"github.com/taikoxyz/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-client/metrics"
	eventiterator "github.com/taikoxyz/taiko-client/pkg/chain_iterator/event_iterator"
	"github.com/taikoxyz/taiko-client/pkg/rpc"
	txListValidator "github.com/taikoxyz/taiko-client/pkg/tx_list_validator"
)

const (
	MaxL1BlocksRead = 1000
)

type L2ChainInserter struct {
	state                         *State            // Driver's state
	rpc                           *rpc.Client       // L1/L2 RPC clients
	throwawayBlocksBuilderPrivKey *ecdsa.PrivateKey // Private key of L2 throwaway blocks builder
	txListValidator               *txListValidator.TxListValidator
}

// NewL2ChainInserter creates a new block inserter instance.
func NewL2ChainInserter(
	ctx context.Context,
	rpc *rpc.Client,
	state *State,
	throwawayBlocksBuilderPrivKey *ecdsa.PrivateKey,
) (*L2ChainInserter, error) {
	return &L2ChainInserter{
		rpc:                           rpc,
		state:                         state,
		throwawayBlocksBuilderPrivKey: throwawayBlocksBuilderPrivKey,
		txListValidator: txListValidator.NewTxListValidator(
			state.maxBlocksGasLimit.Uint64(),
			state.maxBlockNumTxs.Uint64(),
			state.maxTxlistBytes.Uint64(),
			state.minTxGasLimit.Uint64(),
			rpc.L2ChainID,
		),
	}, nil
}

// ProcessL1Blocks fetches all `TaikoL1.BlockProposed` events between given
// L1 block heights, and then tries inserting them into L2 node's block chain.
func (b *L2ChainInserter) ProcessL1Blocks(ctx context.Context, l1End *types.Header) error {
	iter, err := eventiterator.NewBlockProposedIterator(ctx, &eventiterator.BlockProposedIteratorConfig{
		Client:               b.rpc.L1,
		TaikoL1:              b.rpc.TaikoL1,
		StartHeight:          b.state.l1Current.Number,
		EndHeight:            l1End.Number,
		FilterQuery:          nil,
		OnBlockProposedEvent: b.onBlockProposed,
	})
	if err != nil {
		return err
	}

	if err := iter.Iter(); err != nil {
		return err
	}

	b.state.l1Current = l1End
	metrics.DriverL1CurrentHeightGauge.Update(b.state.l1Current.Number.Int64())

	return nil
}

func (b *L2ChainInserter) onBlockProposed(ctx context.Context, event *bindings.TaikoL1ClientBlockProposed) error {
	log.Info(
		"New BlockProposed event",
		"L1Height", event.Raw.BlockNumber,
		"L1Hash", event.Raw.BlockHash,
		"BlockID", event.Id,
	)
	// No need to insert genesis again, its already in L2 block chain.
	if event.Id.Cmp(common.Big0) == 0 {
		return nil
	}

	// Fetch the L2 parent block.
	parent, err := b.rpc.L2ParentByBlockId(ctx, event.Id)
	if err != nil {
		return fmt.Errorf("failed to fetch L2 parent block: %w", err)
	}

	log.Debug("Parent block", "height", parent.Number, "hash", parent.Hash())

	tx, err := b.rpc.L1.TransactionInBlock(
		ctx,
		event.Raw.BlockHash,
		event.Raw.TxIndex,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch original TaikoL1.proposeBlock transaction: %w", err)
	}

	txListBytes, err := encoding.UnpackTxListBytes(tx.Data())
	if err != nil {
		log.Info(
			"Skip the throw away block",
			"blockID", event.Id,
			"hint", "BINARY_NOT_DECODABLE",
			"error", err,
		)
		return nil
	}

	// Check whether the transactions list is valid.
	hint, invalidTxIndex := b.txListValidator.IsTxListValid(event.Id, txListBytes)

	log.Info(
		"Validate transactions list",
		"blockID", event.Id,
		"hint", hint,
		"invalidTxIndex", invalidTxIndex,
	)

	l1Origin := &rawdb.L1Origin{
		BlockID:       event.Id,
		L2BlockHash:   common.Hash{}, // Will be set by taiko-geth.
		L1BlockHeight: new(big.Int).SetUint64(event.Raw.BlockNumber),
		L1BlockHash:   event.Raw.BlockHash,
		Throwaway:     hint != txListValidator.HintOK,
	}

	if event.Meta.Timestamp > uint64(time.Now().Unix()) {
		log.Warn("Future L2 block, waitting", "L2 block timestamp", event.Meta.Timestamp, "now", time.Now().Unix())
		time.Sleep(time.Until(time.Unix(int64(event.Meta.Timestamp), 0)))
	}

	var (
		payloadData  *beacon.ExecutableDataV1
		rpcError     error
		payloadError error
	)
	if hint == txListValidator.HintOK {
		payloadData, rpcError, payloadError = b.insertNewHead(
			ctx,
			event,
			parent,
			b.state.getHeadBlockID(),
			txListBytes,
			l1Origin,
		)
	} else {
		payloadData, rpcError, payloadError = b.insertThrowAwayBlock(
			ctx,
			event,
			parent,
			uint8(hint),
			new(big.Int).SetInt64(int64(invalidTxIndex)),
			b.state.getHeadBlockID(),
			txListBytes,
			l1Origin,
		)
	}

	// RPC errors are recoverable.
	if rpcError != nil {
		return fmt.Errorf("failed to insert new head to L2 node: %w", rpcError)
	}

	if payloadError != nil {
		log.Warn(
			"Ignore invalid block context", "blockID", event.Id, "payloadError", payloadError, "payloadData", payloadData,
		)
		return nil
	}

	log.Debug("Payload data", "payload", payloadData)

	log.Info(
		"🔗 New L2 block inserted",
		"throwaway", l1Origin.Throwaway,
		"blockID", event.Id,
		"height", payloadData.Number,
		"hash", payloadData.BlockHash,
		"lastVerifiedBlockHash", b.state.getLastVerifiedBlockHash(),
		"transactions", len(payloadData.Transactions),
	)

	metrics.DriverL1CurrentHeightGauge.Update(int64(event.Raw.BlockNumber))

	return nil
}

// insertNewHead tries to insert a new head block to the L2 node's local
// block chain through Engine APIs.
func (b *L2ChainInserter) insertNewHead(
	ctx context.Context,
	event *bindings.TaikoL1ClientBlockProposed,
	parent *types.Header,
	headBlockID *big.Int,
	txListBytes []byte,
	l1Origin *rawdb.L1Origin,
) (payloadData *beacon.ExecutableDataV1, rpcError error, payloadError error) {
	log.Debug(
		"Try to insert a new L2 head block",
		"parentNumber", parent.Number,
		"parentHash", parent.Hash(),
		"headBlockID", headBlockID,
		"l1Origin", l1Origin,
	)

	// Insert a TaikoL2.anchor transaction at transactions list head
	var txList []*types.Transaction
	if err := rlp.DecodeBytes(txListBytes, &txList); err != nil {
		log.Info("Ignore invalid txList bytes", "blockID", event.Id)
		return nil, nil, err
	}

	// Assemble a TaikoL2.anchor transaction
	anchorTx, err := b.assembleAnchorTx(
		ctx,
		event.Meta.L1Height,
		event.Meta.L1Hash,
		parent.Number,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create V1TaikoL2.anchor transaction: %w", err)
	}

	txList = append([]*types.Transaction{anchorTx}, txList...)

	if txListBytes, err = rlp.EncodeToBytes(txList); err != nil {
		log.Warn("Encode txList error", "blockID", event.Id, "error", err)
		return nil, nil, err
	}

	payload, rpcErr, payloadErr := b.createExecutionPayloads(
		ctx,
		event,
		parent.Hash(),
		l1Origin,
		headBlockID,
		txListBytes,
	)

	if rpcErr != nil || payloadErr != nil {
		return nil, rpcError, payloadErr
	}

	fc := &beacon.ForkchoiceStateV1{HeadBlockHash: parent.Hash()}

	// Update the fork choice
	fc.HeadBlockHash = payload.BlockHash
	fc.SafeBlockHash = payload.BlockHash
	fcRes, err := b.rpc.L2Engine.ForkchoiceUpdate(ctx, fc, nil)
	if err != nil {
		return nil, err, nil
	}
	if fcRes.PayloadStatus.Status != beacon.VALID {
		return nil, nil, fmt.Errorf("failed to update forkchoice, status: %s", fcRes.PayloadStatus.Status)
	}

	return payload, nil, nil
}

// insertThrowAwayBlock tries to insert a throw away block to the L2 node's local
// block chain through Engine APIs.
func (b *L2ChainInserter) insertThrowAwayBlock(
	ctx context.Context,
	event *bindings.TaikoL1ClientBlockProposed,
	parent *types.Header,
	hint uint8,
	invalidTxIndex *big.Int,
	headBlockID *big.Int,
	txListBytes []byte,
	l1Origin *rawdb.L1Origin,
) (payloadData *beacon.ExecutableDataV1, rpcError error, payloadError error) {
	log.Info(
		"Try to insert a new L2 throwaway block",
		"parentHash", parent.Hash(),
		"headBlockID", headBlockID,
		"l1Origin", l1Origin,
	)

	// Assemble a TaikoL2.invalidateBlock transaction
	opts, err := b.getInvalidateBlockTxOpts(ctx, parent.Number)
	if err != nil {
		return nil, nil, err
	}

	invalidateBlockTx, err := b.rpc.TaikoL2.InvalidateBlock(
		opts,
		txListBytes,
		hint,
		invalidTxIndex,
	)
	if err != nil {
		return nil, nil, err
	}

	throwawayBlockTxListBytes, err := rlp.EncodeToBytes(
		types.Transactions{invalidateBlockTx},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode TaikoL2.InvalidateBlock transaction bytes, err: %w", err)
	}

	return b.createExecutionPayloads(
		ctx,
		event,
		parent.Hash(),
		l1Origin,
		headBlockID,
		throwawayBlockTxListBytes,
	)
}

// createExecutionPayloads creates a new execution payloads through
// Engine APIs.
func (b *L2ChainInserter) createExecutionPayloads(
	ctx context.Context,
	event *bindings.TaikoL1ClientBlockProposed,
	parentHash common.Hash,
	l1Origin *rawdb.L1Origin,
	headBlockID *big.Int,
	txListBytes []byte,
) (payloadData *beacon.ExecutableDataV1, rpcError error, payloadError error) {
	fc := &beacon.ForkchoiceStateV1{HeadBlockHash: parentHash}
	attributes := &beacon.PayloadAttributesV1{
		Timestamp:             event.Meta.Timestamp,
		Random:                event.Meta.MixHash,
		SuggestedFeeRecipient: event.Meta.Beneficiary,
		BlockMetadata: &beacon.BlockMetadata{
			HighestBlockID: headBlockID,
			Beneficiary:    event.Meta.Beneficiary,
			GasLimit:       event.Meta.GasLimit + b.state.anchorTxGasLimit.Uint64(),
			Timestamp:      event.Meta.Timestamp,
			TxList:         txListBytes,
			MixHash:        event.Meta.MixHash,
			ExtraData:      event.Meta.ExtraData,
		},
		L1Origin: l1Origin,
	}

	// TODO: handle payload error more precisely
	// Step 1, prepare a payload
	fcRes, err := b.rpc.L2Engine.ForkchoiceUpdate(ctx, fc, attributes)
	if err != nil {
		return nil, err, nil
	}
	if fcRes.PayloadStatus.Status != beacon.VALID {
		return nil, nil, fmt.Errorf("failed to update forkchoice, status: %s", fcRes.PayloadStatus.Status)
	}
	if fcRes.PayloadID == nil {
		return nil, nil, errors.New("empty payload ID")
	}

	// Step 2, get the payload
	payload, err := b.rpc.L2Engine.GetPayload(ctx, fcRes.PayloadID)
	if err != nil {
		return nil, err, nil
	}

	// Step 3, execute the payload
	execStatus, err := b.rpc.L2Engine.NewPayload(ctx, payload)
	if err != nil {
		return nil, err, nil
	}
	if execStatus.Status != beacon.VALID {
		return nil, nil, fmt.Errorf("failed to execute the newly built block, status: %s", execStatus.Status)
	}

	return payload, nil, nil
}

// getInvalidateBlockTxOpts signs the transaction with a the
// throwaway blocks builder private key.
func (b *L2ChainInserter) getInvalidateBlockTxOpts(ctx context.Context, height *big.Int) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(
		b.throwawayBlocksBuilderPrivKey,
		b.rpc.L2ChainID,
	)
	if err != nil {
		return nil, err
	}

	nonce, err := b.rpc.L2AccountNonce(
		ctx,
		crypto.PubkeyToAddress(b.throwawayBlocksBuilderPrivKey.PublicKey),
		height,
	)
	if err != nil {
		return nil, err
	}

	opts.Nonce = new(big.Int).SetUint64(nonce)
	opts.NoSend = true

	return opts, nil
}
