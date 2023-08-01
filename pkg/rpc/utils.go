package rpc

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/taikoxyz/taiko-client/bindings"
	"github.com/taikoxyz/taiko-client/bindings/encoding"
)

var (
	waitReceiptPollingInterval = 3 * time.Second
	defaultWaitReceiptTimeout  = 1 * time.Minute
	zeroAddress                = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

// GetProtocolStateVariables gets the protocol states from TaikoL1 contract.
func GetProtocolStateVariables(
	taikoL1Client *bindings.TaikoL1Client,
	opts *bind.CallOpts,
) (*bindings.TaikoDataStateVariables, error) {
	stateVars, err := taikoL1Client.GetStateVariables(opts)
	if err != nil {
		return nil, err
	}
	return &stateVars, nil
}

// WaitReceipt keeps waiting until the given transaction has an execution
// receipt to know whether it was reverted or not.
func WaitReceipt(
	ctx context.Context,
	client *EthClient,
	tx *types.Transaction,
) (*types.Receipt, error) {
	ticker := time.NewTicker(waitReceiptPollingInterval)
	defer ticker.Stop()

	var (
		ctxWithTimeout = ctx
		cancel         context.CancelFunc
	)
	if _, ok := ctx.Deadline(); !ok {
		ctxWithTimeout, cancel = context.WithTimeout(ctx, defaultWaitReceiptTimeout)
		defer cancel()
	}

	for {
		select {
		case <-ctxWithTimeout.Done():
			return nil, ctxWithTimeout.Err()
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctxWithTimeout, tx.Hash())
			if err != nil {
				continue
			}

			if receipt.Status != types.ReceiptStatusSuccessful {
				return nil, fmt.Errorf("transaction reverted, hash: %s", tx.Hash())
			}

			return receipt, nil
		}
	}
}

// NeedNewProof checks whether the L2 block still needs a new proof.
func NeedNewProof(
	ctx context.Context,
	cli *Client,
	id *big.Int,
	proverAddress common.Address,
) (bool, error) {
	var parent *types.Header
	if id.Cmp(common.Big1) == 0 {
		header, err := cli.L2.HeaderByNumber(ctx, common.Big0)
		if err != nil {
			return false, err
		}

		parent = header
	} else {
		parentL1Origin, err := cli.WaitL1Origin(ctx, new(big.Int).Sub(id, common.Big1))
		if err != nil {
			return false, err
		}

		if parent, err = cli.L2.HeaderByHash(ctx, parentL1Origin.L2BlockHash); err != nil {
			return false, err
		}
	}

	fc, err := cli.TaikoL1.GetForkChoice(&bind.CallOpts{Context: ctx}, id, parent.Hash(), uint32(parent.GasUsed))
	if err != nil {
		if !strings.Contains(encoding.TryParsingCustomError(err).Error(), "L1_FORK_CHOICE_NOT_FOUND") {
			return false, encoding.TryParsingCustomError(err)
		}

		return true, nil
	}

	if proverAddress == fc.Prover {
		log.Info("📬 Block's proof has already been submitted by current prover", "blockID", id)
		return false, nil
	}

	log.Info(
		"📬 Block's proof has already been submitted by another prover",
		"blockID", id,
		"prover", fc.Prover,
		"provenAt", fc.ProvenAt,
	)

	return false, nil
}

type AccountPoolContent map[string]map[string]*types.Transaction

// ContentFrom fetches a given account's transactions list from a node's transactions pool.
func ContentFrom(
	ctx context.Context,
	rawRPC *rpc.Client,
	address common.Address,
) (AccountPoolContent, error) {
	var result AccountPoolContent
	return result, rawRPC.CallContext(
		ctx,
		&result,
		"txpool_contentFrom",
		address,
	)
}

// GetPendingTxByNonce tries to retrieve a pending transaction with a given nonce in a node's mempool.
func GetPendingTxByNonce(
	ctx context.Context,
	cli *Client,
	address common.Address,
	nonce uint64,
) (*types.Transaction, error) {
	content, err := ContentFrom(ctx, cli.L1RawRPC, address)
	if err != nil {
		return nil, err
	}

	for _, txMap := range content {
		for txNonce, tx := range txMap {
			if txNonce == strconv.Itoa(int(nonce)) {
				return tx, nil
			}
		}
	}

	return nil, nil
}

// SetHead makes a `debug_setHead` RPC call to set the chain's head, should only be used
// for testing purpose.
func SetHead(ctx context.Context, rpc *rpc.Client, headNum *big.Int) error {
	return gethclient.New(rpc).SetHead(ctx, headNum)
}

// StringToBytes32 converts the given string to [32]byte.
func StringToBytes32(str string) [32]byte {
	var b [32]byte
	copy(b[:], []byte(str))

	return b
}

// IsArchiveNode checks if the given node is an archive node.
func IsArchiveNode(ctx context.Context, client *EthClient, l2GenesisHeight uint64) (bool, error) {
	if _, err := client.BalanceAt(ctx, zeroAddress, new(big.Int).SetUint64(l2GenesisHeight)); err != nil {
		if strings.Contains(err.Error(), "missing trie node") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
