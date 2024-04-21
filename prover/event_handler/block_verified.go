package handler

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-client/bindings"
	"github.com/taikoxyz/taiko-client/internal/metrics"
)

// BlockVerifiedEventHandler is responsible for handling the BlockVerified event.
type BlockVerifiedEventHandler struct{}

// Handle handles the BlockVerified event.
func (h *BlockVerifiedEventHandler) Handle(e *bindings.TaikoL1ClientBlockVerified) {
	metrics.ProverLatestVerifiedIDGauge.Set(float64(e.BlockId.Uint64()))

	log.Info(
		"New verified block",
		"blockID", e.BlockId,
		"hash", common.BytesToHash(e.BlockHash[:]),
		"stateRoot", common.BytesToHash(e.StateRoot[:]),
		"prover", e.Prover,
	)
}
