package producer

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
	"github.com/taikoxyz/taiko-client/bindings"
)

func TestRequestProof(t *testing.T) {
	optimisticProofProducer := &OptimisticProofProducer{}

	resCh := make(chan *ProofWithHeader, 1)

	blockID := common.Big32
	header := &types.Header{
		ParentHash:  randHash(),
		UncleHash:   randHash(),
		Coinbase:    common.HexToAddress("0x0000777735367b36bC9B61C50022d9D0700dB4Ec"),
		Root:        randHash(),
		TxHash:      randHash(),
		ReceiptHash: randHash(),
		Difficulty:  common.Big0,
		Number:      common.Big256,
		GasLimit:    1024,
		GasUsed:     1024,
		Time:        uint64(time.Now().Unix()),
		Extra:       randHash().Bytes(),
		MixDigest:   randHash(),
		Nonce:       types.BlockNonce{},
	}
	require.Nil(t, optimisticProofProducer.RequestProof(
		context.Background(),
		&ProofRequestOptions{},
		blockID,
		&bindings.TaikoDataBlockMetadata{},
		header,
		resCh,
	))

	res := <-resCh
	require.Equal(t, res.BlockID, blockID)
	require.Equal(t, res.Header, header)
	require.NotEmpty(t, res.Proof)
}

func TestProofCancel(t *testing.T) {
	optimisticProofProducer := &OptimisticProofProducer{}

	resCh := make(chan *ProofWithHeader, 1)

	blockID := common.Big32
	header := &types.Header{
		ParentHash:  randHash(),
		UncleHash:   randHash(),
		Coinbase:    common.HexToAddress("0x0000777735367b36bC9B61C50022d9D0700dB4Ec"),
		Root:        randHash(),
		TxHash:      randHash(),
		ReceiptHash: randHash(),
		Difficulty:  common.Big0,
		Number:      common.Big256,
		GasLimit:    1024,
		GasUsed:     1024,
		Time:        uint64(time.Now().Unix()),
		Extra:       randHash().Bytes(),
		MixDigest:   randHash(),
		Nonce:       types.BlockNonce{},
	}
	require.Nil(t, optimisticProofProducer.RequestProof(
		context.Background(),
		&ProofRequestOptions{},
		blockID,
		&bindings.TaikoDataBlockMetadata{},
		header,
		resCh,
	))

	// Cancel the proof request, should return nil
	require.Nil(t, optimisticProofProducer.Cancel(context.Background(), blockID))
}

func randHash() common.Hash {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Crit("Failed to generate random bytes", err)
	}
	return common.BytesToHash(b)
}
