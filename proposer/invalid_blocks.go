package proposer

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	globalEpoch uint64 = 0
)

// proposeInvalidBlocksOp tries to propose invalid blocks to TaikoL1 contract
// every `interval` normal propose operations.
func (p *Proposer) proposeInvalidBlocksOp(ctx context.Context, interval uint64) error {
	globalEpoch += 1

	if globalEpoch%interval != 0 {
		return nil
	}

	log.Info("👻 Propose invalid transactions list bytes", "epoch", globalEpoch)

	if err := p.proposeInvalidTxListBytes(ctx); err != nil {
		return fmt.Errorf("failed to propose invalid transaction list bytes: %w", err)
	}

	log.Info("👻 Propose transactions list including invalid transaction", "epoch", globalEpoch)

	if err := p.proposeTxListIncludingInvalidTx(ctx); err != nil {
		return fmt.Errorf("failed to propose transactions list including invalid transaction: %w", err)
	}

	return nil
}

// proposeInvalidTxListBytes commits and proposes an invalid transaction list
// bytes to TaikoL1 contract.
func (p *Proposer) proposeInvalidTxListBytes(ctx context.Context) error {
	return p.commitAndPropose(
		ctx,
		randomBytes(256),
		uint64(rand.Int63n(int64(p.poolContentSplitter.maxGasPerBlock))),
	)
}

// proposeTxListIncludingInvalidTx commits and proposes a validly encoded
// transaction list which including an invalid transaction.
func (p *Proposer) proposeTxListIncludingInvalidTx(ctx context.Context) error {
	invalidTx, err := p.generateInvalidTransaction(ctx)
	if err != nil {
		return err
	}

	txListBytes, err := rlp.EncodeToBytes(types.Transactions{invalidTx})
	if err != nil {
		return err
	}

	return p.commitAndPropose(ctx, txListBytes, invalidTx.Gas())
}

// generateInvalidTransaction creates a transaction with an invalid nonce to
// current L2 world state.
func (p *Proposer) generateInvalidTransaction(ctx context.Context) (*types.Transaction, error) {
	chainId, err := p.rpc.L2.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(p.l1ProposerPrivKey, chainId)
	if err != nil {
		return nil, err
	}

	nonce, err := p.rpc.L2.PendingNonceAt(ctx, crypto.PubkeyToAddress(p.l1ProposerPrivKey.PublicKey))
	if err != nil {
		return nil, err
	}

	opts.GasLimit = 300000
	opts.NoSend = true
	opts.Nonce = new(big.Int).SetUint64(nonce + 1024)

	return p.rpc.TaikoL2.Anchor(opts, common.Big0, common.BytesToHash(randomBytes(32)))
}

// randomBytes generates a random bytes.
func randomBytes(size int) (b []byte) {
	b = make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		log.Crit("Generate random bytes error", "error", err)
	}
	return
}
