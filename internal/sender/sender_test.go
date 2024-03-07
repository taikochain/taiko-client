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

func (s *SenderTestSuite) TestSendTransaction() {
	var (
		opts   = s.sender.GetOpts()
		client = s.RPCClient.L1
		eg     errgroup.Group
	)
	eg.SetLimit(runtime.NumCPU())
	for i := 0; i < 8; i++ {
		i := i
		eg.Go(func() error {
			to := common.BigToAddress(big.NewInt(int64(i)))
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   client.ChainID,
				To:        &to,
				GasFeeCap: opts.GasFeeCap,
				GasTipCap: opts.GasTipCap,
				Gas:       21000000,
				Value:     big.NewInt(1),
				Data:      nil,
			})

			_, err := s.sender.SendTransaction(tx)
			return err
		})
	}
	s.Nil(eg.Wait())

	for _, confirmCh := range s.sender.TxToConfirmChannels() {
		confirm := <-confirmCh
		s.Nil(confirm.Err)
	}
}

func (s *SenderTestSuite) TestSendRawTransaction() {
	nonce, err := s.RPCClient.L1.NonceAt(context.Background(), s.sender.GetOpts().From, nil)
	s.Nil(err)

	var eg errgroup.Group
	eg.SetLimit(runtime.NumCPU())
	for i := 0; i < 5; i++ {
		i := i
		eg.Go(func() error {
			addr := common.BigToAddress(big.NewInt(int64(i)))
			_, err := s.sender.SendRawTransaction(nonce+uint64(i), &addr, big.NewInt(1), nil)
			return err
		})
	}
	s.Nil(eg.Wait())

	for _, confirmCh := range s.sender.TxToConfirmChannels() {
		confirm := <-confirmCh
		s.Nil(confirm.Err)
	}
}

// Test touch max gas price and replacement.
func (s *SenderTestSuite) TestReplacement() {
	send := s.sender
	client := s.RPCClient.L1
	opts := send.GetOpts()

	// Let max gas price be 2 times of the gas fee cap.
	send.MaxGasFee = opts.GasFeeCap.Uint64() * 2

	nonce, err := client.NonceAt(context.Background(), opts.From, nil)
	s.Nil(err)

	pendingNonce, err := client.PendingNonceAt(context.Background(), opts.From)
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
	rawTx, err := send.GetOpts().Signer(send.GetOpts().From, types.NewTx(baseTx))
	s.Nil(err)
	err = client.SendTransaction(context.Background(), rawTx)
	s.Nil(err)

	// Replace the transaction with a higher nonce.
	_, err = send.SendRawTransaction(nonce, &common.Address{}, big.NewInt(1), nil)
	s.Nil(err)

	time.Sleep(time.Second * 6)
	// Send a transaction with a next nonce and let all the transactions be confirmed.
	_, err = send.SendRawTransaction(nonce-1, &common.Address{}, big.NewInt(1), nil)
	s.Nil(err)

	for _, confirmCh := range send.TxToConfirmChannels() {
		confirm := <-confirmCh
		// Check the replaced transaction's gasFeeTap touch the max gas price.
		if confirm.CurrentTx.Nonce() == nonce {
			s.Equal(send.MaxGasFee, confirm.CurrentTx.GasFeeCap().Uint64())
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
	opts := s.sender.GetOpts()

	nonce, err := client.NonceAt(context.Background(), opts.From, nil)
	s.Nil(err)
	pendingNonce, err := client.PendingNonceAt(context.Background(), opts.From)
	s.Nil(err)
	// Run test only if mempool has no pending transactions.
	if pendingNonce > nonce {
		return
	}

	txID, err := send.SendRawTransaction(nonce-3, &common.Address{}, big.NewInt(1), nil)
	s.Nil(err)
	confirm := <-send.TxToConfirmChannel(txID)
	s.Nil(confirm.Err)
	s.Equal(nonce, confirm.CurrentTx.Nonce())
}

func (s *SenderTestSuite) SetupTest() {
	s.ClientTestSuite.SetupTest()
	s.SetL1Automine(true)

	ctx := context.Background()
	priv, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_PROPOSER_PRIVATE_KEY")))
	s.Nil(err)

	s.sender, err = sender.NewSender(ctx, &sender.Config{
		MaxGasFee:      20000000000,
		GasGrowthRate:  50,
		MaxRetrys:      0,
		GasLimit:       2000000,
		MaxWaitingTime: time.Second * 10,
	}, s.RPCClient.L1, priv)
	s.Nil(err)
}

func (s *SenderTestSuite) TearDownTest() {
	s.SetL1Automine(false)
	s.sender.Close()
	s.ClientTestSuite.TearDownTest()
}

func TestSenderTestSuite(t *testing.T) {
	suite.Run(t, new(SenderTestSuite))
}
