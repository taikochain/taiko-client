package sender_test

import (
	"context"
	"math/big"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"

	"github.com/taikoxyz/taiko-client/internal/sender"
	"github.com/taikoxyz/taiko-client/internal/testutils"
)

type SenderTestSuite struct {
	testutils.ClientTestSuite
	sender *sender.Sender
}

func (s *SenderTestSuite) TestNormalSender() {
	var eg errgroup.Group
	eg.SetLimit(runtime.NumCPU())
	for i := 0; i < 5; i++ {
		i := i
		eg.Go(func() error {
			addr := common.BigToAddress(big.NewInt(int64(i)))
			_, err := s.sender.SendRaw(s.sender.Opts.Nonce.Uint64(), &addr, big.NewInt(1), nil)
			return err
		})
	}
	s.Nil(eg.Wait())

	for _, confirmCh := range s.sender.ConfirmChannels() {
		confirm := <-confirmCh
		s.Nil(confirm.Err)
	}
}

// Test touch max gas price and replacement.
func (s *SenderTestSuite) TestReplacement() {
	send := s.sender
	client := s.RPCClient.L1

	// Let max gas price be 2 times of the gas fee cap.
	send.MaxGasFee = send.Opts.GasFeeCap.Uint64() * 2

	nonce, err := client.NonceAt(context.Background(), send.Opts.From, nil)
	s.Nil(err)

	pendingNonce, err := client.PendingNonceAt(context.Background(), send.Opts.From)
	s.Nil(err)
	// Run test only if mempool has no pending transactions.
	if pendingNonce > nonce {
		return
	}

	nonce++
	baseTx := &types.DynamicFeeTx{
		ChainID:   client.ChainID,
		To:        &common.Address{},
		GasFeeCap: big.NewInt(int64(send.MaxGasFee - 1)),
		GasTipCap: big.NewInt(int64(send.MaxGasFee - 1)),
		Nonce:     nonce,
		Gas:       21000,
		Value:     big.NewInt(1),
		Data:      nil,
	}
	rawTx, err := send.Opts.Signer(send.Opts.From, types.NewTx(baseTx))
	s.Nil(err)
	err = client.SendTransaction(context.Background(), rawTx)
	s.Nil(err)

	// Replace the transaction with a higher nonce.
	_, err = send.SendRaw(nonce, &common.Address{}, big.NewInt(1), nil)
	s.Nil(err)

	time.Sleep(time.Second * 6)
	// Send a transaction with a next nonce and let all the transactions be confirmed.
	_, err = send.SendRaw(nonce-1, &common.Address{}, big.NewInt(1), nil)
	s.Nil(err)

	for _, confirmCh := range send.ConfirmChannels() {
		confirm := <-confirmCh
		// Check the replaced transaction's gasFeeTap touch the max gas price.
		if confirm.Tx.Nonce() == nonce {
			s.Equal(send.MaxGasFee, confirm.Tx.GasFeeCap().Uint64())
		}
		s.Nil(confirm.Err)
	}

	_, err = client.TransactionReceipt(context.Background(), rawTx.Hash())
	s.Equal("not found", err.Error())
}

// Test nonce too low.
func (s *SenderTestSuite) TestNonceTooLow() {
	client := s.RPCClient.L1
	send := s.sender

	nonce, err := client.NonceAt(context.Background(), send.Opts.From, nil)
	s.Nil(err)
	pendingNonce, err := client.PendingNonceAt(context.Background(), send.Opts.From)
	s.Nil(err)
	// Run test only if mempool has no pending transactions.
	if pendingNonce > nonce {
		return
	}

	txID, err := send.SendRaw(nonce-3, &common.Address{}, big.NewInt(1), nil)
	s.Nil(err)
	confirm := <-send.ConfirmChannel(txID)
	s.Nil(confirm.Err)
	s.Equal(nonce, confirm.Tx.Nonce())
}

func (s *SenderTestSuite) SetupTest() {
	s.ClientTestSuite.SetupTest()

	ctx := context.Background()
	priv, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_PROPOSER_PRIVATE_KEY")))
	s.Nil(err)

	s.sender, err = sender.NewSender(ctx, &sender.Config{
		MaxGasFee:      20000000000,
		GasGrowthRate:  50,
		RetryTimes:     0,
		GasLimit:       2000000,
		MaxWaitingTime: time.Second * 10,
	}, s.RPCClient.L1, priv)
	s.Nil(err)
}

func (s *SenderTestSuite) TearDownTest() {
	s.sender.Close()
	s.ClientTestSuite.TearDownTest()
}

func TestSenderTestSuite(t *testing.T) {
	suite.Run(t, new(SenderTestSuite))
}
