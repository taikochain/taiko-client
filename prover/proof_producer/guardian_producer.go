package producer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-client/bindings"
	"github.com/taikoxyz/taiko-client/bindings/encoding"
)

// GuardianProofProducer always returns an optimistic (dummy) proof.
type GuardianProofProducer struct {
	returnLivenessBond bool
	*DummyProofProducer
}

func NewGuardianProofProducer(returnLivenessBond bool) *GuardianProofProducer {
	gp := &GuardianProofProducer{
		returnLivenessBond: returnLivenessBond,
	}

	if !returnLivenessBond {
		gp.DummyProofProducer = new(DummyProofProducer)
	}

	return gp
}

// RequestProof implements the ProofProducer interface.
func (g *GuardianProofProducer) RequestProof(
	ctx context.Context,
	opts *ProofRequestOptions,
	blockID *big.Int,
	meta *bindings.TaikoDataBlockMetadata,
	header *types.Header,
) (*ProofWithHeader, error) {
	log.Info(
		"Request guardian proof",
		"blockID", blockID,
		"coinbase", meta.Coinbase,
		"height", header.Number,
		"hash", header.Hash(),
	)

	if g.returnLivenessBond {
		return &ProofWithHeader{
			BlockID: blockID,
			Meta:    meta,
			Header:  header,
			Proof:   crypto.Keccak256([]byte("RETURN_LIVENESS_BOND")),
			Opts:    opts,
			Tier:    g.Tier(),
		}, nil
	}

	return g.DummyProofProducer.RequestProof(ctx, opts, blockID, meta, header, g.Tier())
}

// Tier implements the ProofProducer interface.
func (g *GuardianProofProducer) Tier() uint16 {
	return encoding.TierGuardianID
}

// Cancellable implements the ProofProducer interface.
func (g *GuardianProofProducer) Cancellable() bool {
	return false
}

// Cancel cancels an existing proof generation.
func (g *GuardianProofProducer) Cancel(ctx context.Context, blockID *big.Int) error {
	return nil
}
