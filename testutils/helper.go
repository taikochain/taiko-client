package testutils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/phayes/freeport"
	"github.com/taikoxyz/taiko-client/bindings"
	"github.com/taikoxyz/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-client/prover/http"
)

func ProposeInvalidTxListBytes(s *ClientTestSuite, proposer Proposer) {
	invalidTxListBytes := RandomBytes(256)

	s.Nil(proposer.ProposeTxList(context.Background(), &encoding.TaikoL1BlockMetadataInput{
		Beneficiary:     proposer.L2SuggestedFeeRecipient(),
		TxListHash:      crypto.Keccak256Hash(invalidTxListBytes),
		TxListByteStart: common.Big0,
		TxListByteEnd:   new(big.Int).SetUint64(uint64(len(invalidTxListBytes))),
		CacheTxListInfo: false,
	}, invalidTxListBytes, 1, nil))
}

func ProposeAndInsertEmptyBlocks(
	s *ClientTestSuite,
	proposer Proposer,
	calldataSyncer CalldataSyncer,
) []*bindings.TaikoL1ClientBlockProposed {
	var events []*bindings.TaikoL1ClientBlockProposed

	l1Head, err := s.RpcClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	sink := make(chan *bindings.TaikoL1ClientBlockProposed)

	sub, err := s.RpcClient.TaikoL1.WatchBlockProposed(nil, sink, nil, nil)
	s.Nil(err)
	defer func() {
		sub.Unsubscribe()
		close(sink)
	}()

	// RLP encoded empty list
	var emptyTxs []types.Transaction
	encoded, err := rlp.EncodeToBytes(emptyTxs)
	s.Nil(err)

	s.Nil(proposer.ProposeTxList(context.Background(), &encoding.TaikoL1BlockMetadataInput{
		Beneficiary:     proposer.L2SuggestedFeeRecipient(),
		TxListHash:      crypto.Keccak256Hash(encoded),
		TxListByteStart: common.Big0,
		TxListByteEnd:   new(big.Int).SetUint64(uint64(len(encoded))),
		CacheTxListInfo: false,
	}, encoded, 0, nil))

	ProposeInvalidTxListBytes(s, proposer)

	// Zero byte txList
	s.Nil(proposer.ProposeEmptyBlockOp(context.Background()))

	events = append(events, []*bindings.TaikoL1ClientBlockProposed{<-sink, <-sink, <-sink}...)

	_, isPending, err := s.RpcClient.L1.TransactionByHash(context.Background(), events[len(events)-1].Raw.TxHash)
	s.Nil(err)
	s.False(isPending)

	newL1Head, err := s.RpcClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)
	s.Greater(newL1Head.Number.Uint64(), l1Head.Number.Uint64())

	syncProgress, err := s.RpcClient.L2.SyncProgress(context.Background())
	s.Nil(err)
	s.Nil(syncProgress)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.Nil(calldataSyncer.ProcessL1Blocks(ctx, newL1Head))

	return events
}

// ProposeAndInsertValidBlock proposes an valid tx list and then insert it
// into L2 execution engine's local chain.
func ProposeAndInsertValidBlock(
	s *ClientTestSuite,
	proposer Proposer,
	calldataSyncer CalldataSyncer,
) *bindings.TaikoL1ClientBlockProposed {
	l1Head, err := s.RpcClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	l2Head, err := s.RpcClient.L2.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	// Propose txs in L2 execution engine's mempool
	sink := make(chan *bindings.TaikoL1ClientBlockProposed)

	sub, err := s.RpcClient.TaikoL1.WatchBlockProposed(nil, sink, nil, nil)
	s.Nil(err)
	defer func() {
		sub.Unsubscribe()
		close(sink)
	}()

	baseFee, err := s.RpcClient.TaikoL2.GetBasefee(nil, 0, uint32(l2Head.GasUsed))
	s.Nil(err)

	nonce, err := s.RpcClient.L2.PendingNonceAt(context.Background(), s.TestAddr)
	s.Nil(err)

	tx := types.NewTransaction(
		nonce,
		common.BytesToAddress(RandomBytes(32)),
		common.Big1,
		100000,
		baseFee,
		[]byte{},
	)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(s.RpcClient.L2ChainID), s.TestAddrPrivKey)
	s.Nil(err)
	s.Nil(s.RpcClient.L2.SendTransaction(context.Background(), signedTx))

	s.Nil(proposer.ProposeOp(context.Background()))

	event := <-sink

	_, isPending, err := s.RpcClient.L1.TransactionByHash(context.Background(), event.Raw.TxHash)
	s.Nil(err)
	s.False(isPending)

	receipt, err := s.RpcClient.L1.TransactionReceipt(context.Background(), event.Raw.TxHash)
	s.Nil(err)
	s.Equal(types.ReceiptStatusSuccessful, receipt.Status)

	newL1Head, err := s.RpcClient.L1.HeaderByNumber(context.Background(), nil)
	s.Nil(err)
	s.Greater(newL1Head.Number.Uint64(), l1Head.Number.Uint64())

	syncProgress, err := s.RpcClient.L2.SyncProgress(context.Background())
	s.Nil(err)
	s.Nil(syncProgress)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.Nil(calldataSyncer.ProcessL1Blocks(ctx, newL1Head))

	_, err = s.RpcClient.L2.HeaderByNumber(context.Background(), nil)
	s.Nil(err)

	return event
}

func DepositEtherToL2(s *ClientTestSuite, depositerPrivKey *ecdsa.PrivateKey, recipient common.Address) {
	config, err := s.RpcClient.TaikoL1.GetConfig(nil)
	s.Nil(err)

	opts, err := bind.NewKeyedTransactorWithChainID(depositerPrivKey, s.RpcClient.L1ChainID)
	s.Nil(err)
	opts.Value = config.EthDepositMinAmount

	for i := 0; i < int(config.EthDepositMinCountPerBlock); i++ {
		_, err = s.RpcClient.TaikoL1.DepositEtherToL2(opts, recipient)
		s.Nil(err)
	}
}

// HTTPServer starts a new prover server that has channel listeners to respond and react
// to requests for capacity, which provers can call.
func HTTPServer(s *ClientTestSuite, port int) (*http.Server, func(), error) {
	l1ProverPrivKey, err := crypto.ToECDSA(common.Hex2Bytes(os.Getenv("L1_PROVER_PRIVATE_KEY")))
	s.Nil(err)

	serverOpts := http.NewServerOpts{
		ProverPrivateKey:         l1ProverPrivKey,
		MinProofFee:              big.NewInt(1),
		MaxCapacity:              10,
		RequestCurrentCapacityCh: make(chan struct{}),
		ReceiveCurrentCapacityCh: make(chan uint64),
	}

	srv, err := http.NewServer(serverOpts)
	s.Nil(err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-serverOpts.RequestCurrentCapacityCh:
				serverOpts.ReceiveCurrentCapacityCh <- 100
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		_ = srv.Start(fmt.Sprintf(":%v", port))
	}()

	return srv, func() {
		cancel()
		_ = srv.Shutdown(ctx)
	}, err
}

// RandomHash generates a random blob of data and returns it as a hash.
func RandomHash() common.Hash {
	var hash common.Hash
	if n, err := rand.Read(hash[:]); n != common.HashLength || err != nil {
		panic(err)
	}
	return hash
}

// RandomBytes generates a random bytes.
func RandomBytes(size int) (b []byte) {
	b = make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		log.Crit("Generate random bytes error", "error", err)
	}
	return
}

// RandomPort returns a local free random port.
func RandomPort() int {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Crit("Failed to get local free random port", "err", err)
	}
	return port
}

// SignatureFromRSV creates the signature bytes from r,s,v.
func SignatureFromRSV(r, s string, v byte) []byte {
	return append(append(hexutil.MustDecode(r), hexutil.MustDecode(s)...), v)
}
