package producer

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/taikoxyz/taiko-client/bindings"
)

func TestDummyProducerRequestProof(t *testing.T) {
	producer := &DummyProofProducer{}

	resCh := make(chan *ProofWithHeader, 1)

	var tier uint16 = 1024

	blockID := common.Big32
	header := &types.Header{
		ParentHash:  randHash(),
		UncleHash:   randHash(),
		Coinbase:    common.BytesToAddress(randHash().Bytes()),
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
	require.Nil(t, producer.RequestProof(
		context.Background(),
		&ProofRequestOptions{},
		blockID,
		&bindings.TaikoDataBlockMetadata{},
		header,
		tier,
		resCh,
	))

	res := <-resCh
	require.Equal(t, res.BlockID, blockID)
	require.Equal(t, res.Header, header)
	require.Equal(t, tier, res.Tier)
	require.NotEmpty(t, res.Proof)
}
