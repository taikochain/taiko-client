package prover

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// validateAnchorTx checks whether the given transaction is a successfully
// executed `TaikoL2.anchor` transaction.
func (p *Prover) validateAnchorTx(ctx context.Context, tx *types.Transaction) error {
	if tx.To() == nil || *tx.To() != p.cfg.TaikoL2Address {
		return fmt.Errorf("invalid TaikoL2.anchor to: %s", tx.To())
	}

	// TODO: verify sender

	method, err := p.taikoL2Abi.MethodById(tx.Data())
	if err != nil || method.Name != "anchor" {
		return fmt.Errorf("invalid method method, err: %w, methodName: %s", err, method.Name)
	}

	receipt, err := p.l2RPC.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return fmt.Errorf("failed to get TaikoL2.anchor receipt, err: %w", err)
	}

	log.Info("TaikoL2.anchor transaction found", "height", receipt.BlockNumber, "status", receipt.Status)

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("invalid TaikoL2.anchor receipt status: %d", receipt.Status)
	}

	if len(receipt.Logs) == 0 {
		return fmt.Errorf("HeaderExchanged event not found in TaikoL2.anchor receipt")
	}

	return nil
}
