package calldata

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"

	"github.com/taikoxyz/taiko-client/bindings"
	"github.com/taikoxyz/taiko-client/driver/chain_syncer/beaconsync"
	"github.com/taikoxyz/taiko-client/driver/state"
	"github.com/taikoxyz/taiko-client/internal/testutils"
	"github.com/taikoxyz/taiko-client/internal/utils"
	"github.com/taikoxyz/taiko-client/pkg/rpc"
	"github.com/taikoxyz/taiko-client/proposer"
)

type CalldataSyncerTestSuite struct {
	testutils.ClientTestSuite
	s *Syncer
	p testutils.Proposer
}

func (s *CalldataSyncerTestSuite) SetupTest() {
	s.ClientTestSuite.SetupTest()

	state2, err := state.New(context.Background(), s.RPCClient)
	s.Nil(err)

	syncer, err := NewSyncer(
		context.Background(),
		s.RPCClient,
		state2,
		beaconsync.NewSyncProgressTracker(s.RPCClient.L2, 1*time.Hour),
		0,
	)
	s.Nil(err)
	s.s = syncer

	s.initProposer()
}
func (s *CalldataSyncerTestSuite) TestCancelNewSyncer() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	syncer, err := NewSyncer(
		ctx,
		s.RPCClient,
		s.s.state,
		s.s.progressTracker,
		0,
	)
	s.Nil(syncer)
	s.NotNil(err)
}

func (s *CalldataSyncerTestSuite) TestProcessL1Blocks() {
	s.Nil(s.s.ProcessL1Blocks(context.Background()))
}

func (s *CalldataSyncerTestSuite) TestProcessL1BlocksReorg() {
	s.ProposeAndInsertEmptyBlocks(s.p, s.s)
	s.Nil(s.s.ProcessL1Blocks(context.Background()))
}

func (s *CalldataSyncerTestSuite) TestOnBlockProposed() {
	s.Nil(s.s.onBlockProposed(
		context.Background(),
		&bindings.TaikoL1ClientBlockProposed{BlockId: common.Big0},
		func() {},
	))
	s.NotNil(s.s.onBlockProposed(
		context.Background(),
		&bindings.TaikoL1ClientBlockProposed{BlockId: common.Big1},
		func() {},
	))
}

func (s *CalldataSyncerTestSuite) TestInsertNewHead() {
	parent, err := s.s.rpc.L2.HeaderByNumber(context.Background(), nil)
	s.Nil(err)
	l1Head, err := s.s.rpc.L1.BlockByNumber(context.Background(), nil)
	s.Nil(err)
	_, err = s.s.insertNewHead(
		context.Background(),
		&bindings.TaikoL1ClientBlockProposed{
			BlockId: common.Big1,
			Meta: bindings.TaikoDataBlockMetadata{
				Id:         1,
				L1Height:   l1Head.NumberU64(),
				L1Hash:     l1Head.Hash(),
				Coinbase:   common.BytesToAddress(testutils.RandomBytes(1024)),
				BlobHash:   testutils.RandomHash(),
				Difficulty: testutils.RandomHash(),
				GasLimit:   utils.RandUint32(nil),
				Timestamp:  uint64(time.Now().Unix()),
			},
		},
		parent,
		common.Big2,
		[]byte{},
		&rawdb.L1Origin{
			BlockID:       common.Big1,
			L1BlockHeight: common.Big1,
			L1BlockHash:   testutils.RandomHash(),
		},
	)
	s.Nil(err)
}

func (s *CalldataSyncerTestSuite) TestTreasuryIncomeAllAnchors() {
	treasury := common.HexToAddress(os.Getenv("TREASURY"))
	s.NotZero(treasury.Big().Uint64())

	balance, err := s.RPCClient.L2.BalanceAt(context.Background(), treasury, nil)
	s.Nil(err)

	headBefore, err := s.RPCClient.L2.BlockNumber(context.Background())
	s.Nil(err)

	s.ProposeAndInsertEmptyBlocks(s.p, s.s)

	headAfter, err := s.RPCClient.L2.BlockNumber(context.Background())
	s.Nil(err)

	balanceAfter, err := s.RPCClient.L2.BalanceAt(context.Background(), treasury, nil)
	s.Nil(err)

	s.Greater(headAfter, headBefore)
	s.Zero(balanceAfter.Cmp(balance))
}

func (s *CalldataSyncerTestSuite) TestTreasuryIncome() {
	treasury := common.HexToAddress(os.Getenv("TREASURY"))
	s.NotZero(treasury.Big().Uint64())

	balance, err := s.RPCClient.L2.BalanceAt(context.Background(), treasury, nil)
	s.Nil(err)

	headBefore, err := s.RPCClient.L2.BlockNumber(context.Background())
	s.Nil(err)

	s.ProposeAndInsertEmptyBlocks(s.p, s.s)
	s.ProposeAndInsertValidBlock(s.p, s.s)

	headAfter, err := s.RPCClient.L2.BlockNumber(context.Background())
	s.Nil(err)

	balanceAfter, err := s.RPCClient.L2.BalanceAt(context.Background(), treasury, nil)
	s.Nil(err)

	s.Greater(headAfter, headBefore)
	s.True(balanceAfter.Cmp(balance) > 0)

	var hasNoneAnchorTxs bool
	for i := headBefore + 1; i <= headAfter; i++ {
		block, err := s.RPCClient.L2.BlockByNumber(context.Background(), new(big.Int).SetUint64(i))
		s.Nil(err)
		s.GreaterOrEqual(block.Transactions().Len(), 1)
		s.Greater(block.BaseFee().Uint64(), uint64(0))

		for j, tx := range block.Transactions() {
			if j == 0 {
				continue
			}

			hasNoneAnchorTxs = true
			receipt, err := s.RPCClient.L2.TransactionReceipt(context.Background(), tx.Hash())
			s.Nil(err)

			fee := new(big.Int).Mul(block.BaseFee(), new(big.Int).SetUint64(receipt.GasUsed))

			balance = new(big.Int).Add(balance, fee)
		}
	}

	s.True(hasNoneAnchorTxs)
	s.Zero(balanceAfter.Cmp(balance))
}

func (s *CalldataSyncerTestSuite) initProposer() {
	prop := new(proposer.Proposer)
	l1ProposerPrivKey, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_PROPOSER_PRIVATE_KEY")))
	s.Nil(err)

	s.Nil(prop.InitFromConfig(context.Background(), &proposer.Config{
		ClientConfig: &rpc.ClientConfig{
			L1Endpoint:        os.Getenv("L1_NODE_WS_ENDPOINT"),
			L2Endpoint:        os.Getenv("L2_EXECUTION_ENGINE_WS_ENDPOINT"),
			TaikoL1Address:    common.HexToAddress(os.Getenv("TAIKO_L1_ADDRESS")),
			TaikoL2Address:    common.HexToAddress(os.Getenv("TAIKO_L2_ADDRESS")),
			TaikoTokenAddress: common.HexToAddress(os.Getenv("TAIKO_TOKEN_ADDRESS")),
		},
		AssignmentHookAddress:      common.HexToAddress(os.Getenv("ASSIGNMENT_HOOK_ADDRESS")),
		L1ProposerPrivKey:          l1ProposerPrivKey,
		L2SuggestedFeeRecipient:    common.HexToAddress(os.Getenv("L2_SUGGESTED_FEE_RECIPIENT")),
		ProposeInterval:            1024 * time.Hour,
		MaxProposedTxListsPerEpoch: 1,
		WaitReceiptTimeout:         12 * time.Second,
		ProverEndpoints:            s.ProverEndpoints,
		OptimisticTierFee:          common.Big256,
		SgxTierFee:                 common.Big256,
		MaxTierFeePriceBumps:       3,
		TierFeePriceBump:           common.Big2,
		L1BlockBuilderTip:          common.Big0,
		TxmgrConfigs: &txmgr.CLIConfig{
			L1RPCURL:                  os.Getenv("L1_NODE_WS_ENDPOINT"),
			NumConfirmations:          0,
			SafeAbortNonceTooLowCount: txmgr.DefaultBatcherFlagValues.SafeAbortNonceTooLowCount,
			PrivateKey:                common.Bytes2Hex(crypto.FromECDSA(l1ProposerPrivKey)),
			FeeLimitMultiplier:        txmgr.DefaultBatcherFlagValues.FeeLimitMultiplier,
			FeeLimitThresholdGwei:     txmgr.DefaultBatcherFlagValues.FeeLimitThresholdGwei,
			MinBaseFeeGwei:            txmgr.DefaultBatcherFlagValues.MinBaseFeeGwei,
			MinTipCapGwei:             txmgr.DefaultBatcherFlagValues.MinTipCapGwei,
			ResubmissionTimeout:       txmgr.DefaultBatcherFlagValues.ResubmissionTimeout,
			ReceiptQueryInterval:      1 * time.Second,
			NetworkTimeout:            txmgr.DefaultBatcherFlagValues.NetworkTimeout,
			TxSendTimeout:             txmgr.DefaultBatcherFlagValues.TxSendTimeout,
			TxNotInMempoolTimeout:     txmgr.DefaultBatcherFlagValues.TxNotInMempoolTimeout,
		},
	}))

	s.p = prop
}

func TestCalldataSyncerTestSuite(t *testing.T) {
	suite.Run(t, new(CalldataSyncerTestSuite))
}
