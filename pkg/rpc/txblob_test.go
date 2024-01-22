package rpc

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestBlockTx(t *testing.T) {
	l1Client, err := NewEthClient(nil, os.Getenv("L1_NODE_WS_ENDPOINT"), time.Second*20)
	assert.NoError(t, err)

	sk, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_PROPOSER_PRIVATE_KEY")))
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chainID, err := l1Client.ChainID(ctx)
	assert.NoError(t, err)

	opts, err := bind.NewKeyedTransactorWithChainID(sk, chainID)
	assert.NoError(t, err)
	opts.Context = ctx

	tx, err := l1Client.TransactBlobTx(opts, []byte("ssssssssss"))
	assert.NoError(t, err)

	receipt, err := bind.WaitMined(ctx, l1Client, tx)
	assert.NoError(t, err)
	t.Log(receipt)

}
