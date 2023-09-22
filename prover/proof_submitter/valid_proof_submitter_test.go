package submitter

import (
	"bytes"
	"context"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
	"github.com/taikoxyz/taiko-client/bindings"
	"github.com/taikoxyz/taiko-client/driver/chain_syncer/beaconsync"
	"github.com/taikoxyz/taiko-client/driver/chain_syncer/calldata"
	"github.com/taikoxyz/taiko-client/driver/state"
	"github.com/taikoxyz/taiko-client/proposer"
	proofProducer "github.com/taikoxyz/taiko-client/prover/proof_producer"
	"github.com/taikoxyz/taiko-client/testutils"
)

type ProofSubmitterTestSuite struct {
	testutils.ClientTestSuite
	validProofSubmitter *ValidProofSubmitter
	calldataSyncer      *calldata.Syncer
	proposer            *proposer.Proposer
	validProofCh        chan *proofProducer.ProofWithHeader
	invalidProofCh      chan *proofProducer.ProofWithHeader
}

func (s *ProofSubmitterTestSuite) SetupTest() {
	s.ClientTestSuite.SetupTest()

	l1ProverPrivKey, err := crypto.ToECDSA(common.Hex2Bytes(os.Getenv("L1_PROVER_PRIVATE_KEY")))
	s.NoError(err)

	s.validProofCh = make(chan *proofProducer.ProofWithHeader, 1024)
	s.invalidProofCh = make(chan *proofProducer.ProofWithHeader, 1024)

	s.validProofSubmitter, err = NewValidProofSubmitter(
		s.RpcClient,
		&proofProducer.DummyProofProducer{},
		s.validProofCh,
		common.HexToAddress(os.Getenv("TAIKO_L2_ADDRESS")),
		l1ProverPrivKey,
		&sync.Mutex{},
		false,
		"test",
		1,
		12*time.Second,
		10*time.Second,
		nil,
	)
	s.NoError(err)
	ctx := context.Background()
	// Init calldata syncer
	testState, err := state.New(ctx, s.RpcClient)
	s.NoError(err)

	tracker := beaconsync.NewSyncProgressTracker(s.RpcClient.L2, 30*time.Second)

	s.calldataSyncer, err = calldata.NewSyncer(
		ctx,
		s.RpcClient,
		testState,
		tracker,
		common.HexToAddress(os.Getenv("L1_SIGNAL_SERVICE_CONTRACT_ADDRESS")),
	)
	s.NoError(err)

	// Init proposer
	l1ProposerPrivKey, err := crypto.ToECDSA(common.Hex2Bytes(os.Getenv("L1_PROPOSER_PRIVATE_KEY")))
	s.NoError(err)
	proposeInterval := 1024 * time.Hour // No need to periodically propose transactions list in unit tests

	cfg := &proposer.Config{
		L1Endpoint:                         os.Getenv("L1_NODE_WS_ENDPOINT"),
		L2Endpoint:                         os.Getenv("L2_EXECUTION_ENGINE_WS_ENDPOINT"),
		TaikoL1Address:                     common.HexToAddress(os.Getenv("TAIKO_L1_ADDRESS")),
		TaikoL2Address:                     common.HexToAddress(os.Getenv("TAIKO_L2_ADDRESS")),
		TaikoTokenAddress:                  common.HexToAddress(os.Getenv("TAIKO_TOKEN_ADDRESS")),
		L1ProposerPrivKey:                  l1ProposerPrivKey,
		L2SuggestedFeeRecipient:            common.HexToAddress(os.Getenv("L2_SUGGESTED_FEE_RECIPIENT")),
		ProposeInterval:                    &proposeInterval,
		MaxProposedTxListsPerEpoch:         1,
		WaitReceiptTimeout:                 10 * time.Second,
		ProverEndpoints:                    s.ProverEndpoints,
		BlockProposalFee:                   big.NewInt(1000),
		BlockProposalFeeIterations:         3,
		BlockProposalFeeIncreasePercentage: common.Big2,
	}
	ep, err := proposer.EndpointFromConfig(ctx, cfg)
	s.NoError(err)
	s.proposer, err = proposer.New(ctx, ep, cfg)
	s.NoError(err)
}

func (s *ProofSubmitterTestSuite) TestValidProofSubmitterRequestProofDeadlineExceeded() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s.ErrorContains(
		s.validProofSubmitter.RequestProof(
			ctx, &bindings.TaikoL1ClientBlockProposed{BlockId: common.Big256}), "context deadline exceeded",
	)
}

func (s *ProofSubmitterTestSuite) TestValidProofSubmitterSubmitProofMetadataNotFound() {
	s.Error(
		s.validProofSubmitter.SubmitProof(
			context.Background(), &proofProducer.ProofWithHeader{
				BlockID: common.Big256,
				Meta:    &bindings.TaikoDataBlockMetadata{},
				Header:  &types.Header{},
				ZkProof: bytes.Repeat([]byte{0xff}, 100),
			},
		),
	)
}

func (s *ProofSubmitterTestSuite) TestValidSubmitProofs() {
	events := testutils.ProposeAndInsertEmptyBlocks(&s.ClientTestSuite, s.proposer, s.calldataSyncer)

	for _, e := range events {
		s.Nil(s.validProofSubmitter.RequestProof(context.Background(), e))
		proofWithHeader := <-s.validProofCh
		s.Nil(s.validProofSubmitter.SubmitProof(context.Background(), proofWithHeader))
	}
}

func (s *ProofSubmitterTestSuite) TestValidProofSubmitterRequestProofCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.AfterFunc(2*time.Second, func() {
			cancel()
		})
	}()

	s.ErrorContains(
		s.validProofSubmitter.RequestProof(
			ctx, &bindings.TaikoL1ClientBlockProposed{BlockId: common.Big256}), "context canceled",
	)
}

func TestProofSubmitterTestSuite(t *testing.T) {
	suite.Run(t, new(ProofSubmitterTestSuite))
}
