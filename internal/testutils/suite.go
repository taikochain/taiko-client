package testutils

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"net/url"
	"os"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/suite"

	"github.com/taikoxyz/taiko-client/bindings"
	"github.com/taikoxyz/taiko-client/internal/utils"
	"github.com/taikoxyz/taiko-client/pkg/jwt"
	"github.com/taikoxyz/taiko-client/pkg/rpc"
	"github.com/taikoxyz/taiko-client/prover/server"
)

type ClientTestSuite struct {
	suite.Suite
	testnetL1SnapshotID string
	RPCClient           *rpc.Client
	TestAddrPrivKey     *ecdsa.PrivateKey
	TestAddr            common.Address
	ProverEndpoints     []*url.URL
	AddressManager      *bindings.AddressManager
	proverServer        *server.ProverServer
}

func (s *ClientTestSuite) SetupTest() {
	utils.LoadEnv()
	// Default logger
	glogger := log.NewGlogHandler(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelInfo, true))
	log.SetDefault(log.NewLogger(glogger))

	testAddrPrivKey, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_PROPOSER_PRIVATE_KEY")))
	s.Nil(err)

	s.TestAddrPrivKey = testAddrPrivKey
	s.TestAddr = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

	jwtSecret, err := jwt.ParseSecretFromFile(os.Getenv("JWT_SECRET"))
	s.Nil(err)
	s.NotEmpty(jwtSecret)

	rpcCli, err := rpc.NewClient(context.Background(), &rpc.ClientConfig{
		L1Endpoint:            os.Getenv("L1_NODE_WS_ENDPOINT"),
		L2Endpoint:            os.Getenv("L2_EXECUTION_ENGINE_WS_ENDPOINT"),
		TaikoL1Address:        common.HexToAddress(os.Getenv("TAIKO_L1_ADDRESS")),
		TaikoL2Address:        common.HexToAddress(os.Getenv("TAIKO_L2_ADDRESS")),
		TaikoTokenAddress:     common.HexToAddress(os.Getenv("TAIKO_TOKEN_ADDRESS")),
		GuardianProverAddress: common.HexToAddress(os.Getenv("GUARDIAN_PROVER_CONTRACT_ADDRESS")),
		L2EngineEndpoint:      os.Getenv("L2_EXECUTION_ENGINE_AUTH_ENDPOINT"),
		JwtSecret:             string(jwtSecret),
		RetryInterval:         backoff.DefaultMaxInterval,
	})
	s.Nil(err)
	s.RPCClient = rpcCli

	l1ProverPrivKey, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_PROVER_PRIVATE_KEY")))
	s.Nil(err)

	s.ProverEndpoints = []*url.URL{LocalRandomProverEndpoint()}
	s.proverServer = s.NewTestProverServer(l1ProverPrivKey, s.ProverEndpoints[0])

	balance, err := rpcCli.TaikoToken.BalanceOf(nil, crypto.PubkeyToAddress(l1ProverPrivKey.PublicKey))
	s.Nil(err)

	if balance.Cmp(common.Big0) == 0 {
		ownerPrivKey, err := crypto.ToECDSA(common.FromHex(os.Getenv("L1_CONTRACT_OWNER_PRIVATE_KEY")))
		s.Nil(err)

		// Transfer some tokens to provers.
		balance, err := rpcCli.TaikoToken.BalanceOf(nil, crypto.PubkeyToAddress(ownerPrivKey.PublicKey))
		s.Nil(err)
		s.Greater(balance.Cmp(common.Big0), 0)

		opts, err := bind.NewKeyedTransactorWithChainID(ownerPrivKey, rpcCli.L1.ChainID)
		s.Nil(err)
		proverBalance := new(big.Int).Div(balance, common.Big2)
		s.Greater(proverBalance.Cmp(common.Big0), 0)

		tx, err := rpcCli.TaikoToken.Transfer(opts, crypto.PubkeyToAddress(l1ProverPrivKey.PublicKey), proverBalance)
		s.Nil(err)
		_, err = rpc.WaitReceipt(context.Background(), rpcCli.L1, tx)
		s.Nil(err)

		// Increase allowance for AssignmentHook and TaikoL1
		s.setAllowance(l1ProverPrivKey)
		s.setAllowance(ownerPrivKey)
	}
	s.testnetL1SnapshotID = s.SetL1Snapshot()
}

func (s *ClientTestSuite) setAllowance(key *ecdsa.PrivateKey) {
	decimal, err := s.RPCClient.TaikoToken.Decimals(nil)
	s.Nil(err)

	bigInt := new(big.Int).Exp(big.NewInt(1_000_000_000), new(big.Int).SetUint64(uint64(decimal)), nil)

	opts, err := bind.NewKeyedTransactorWithChainID(key, s.RPCClient.L1.ChainID)
	s.Nil(err)

	_, err = s.RPCClient.TaikoToken.Approve(
		opts,
		common.HexToAddress(os.Getenv("ASSIGNMENT_HOOK_ADDRESS")),
		bigInt,
	)
	s.Nil(err)

	tx, err := s.RPCClient.TaikoToken.Approve(
		opts,
		common.HexToAddress(os.Getenv("TAIKO_L1_ADDRESS")),
		bigInt,
	)
	s.Nil(err)

	_, err = rpc.WaitReceipt(context.Background(), s.RPCClient.L1, tx)
	s.Nil(err)
}

func (s *ClientTestSuite) TearDownTest() {
	s.RevertL1Snapshot(s.testnetL1SnapshotID)

	s.Nil(rpc.SetHead(context.Background(), s.RPCClient.L2, common.Big0))
	s.Nil(s.proverServer.Shutdown(context.Background()))
}

func (s *ClientTestSuite) SetL1Automine(automine bool) {
	s.Nil(s.RPCClient.L1.CallContext(context.Background(), nil, "evm_setAutomine", automine))
}

func (s *ClientTestSuite) IncreaseTime(time uint64) {
	var result uint64
	s.Nil(s.RPCClient.L1.CallContext(context.Background(), &result, "evm_increaseTime", time))
	s.NotNil(result)
}

func (s *ClientTestSuite) SetL1Snapshot() string {
	var snapshotID string
	s.Nil(s.RPCClient.L1.CallContext(context.Background(), &snapshotID, "evm_snapshot"))
	s.NotEmpty(snapshotID)
	return snapshotID
}

func (s *ClientTestSuite) RevertL1Snapshot(snapshotID string) {
	var revertRes bool
	s.Nil(s.RPCClient.L1.CallContext(context.Background(), &revertRes, "evm_revert", snapshotID))
	s.True(revertRes)
}

func (s *ClientTestSuite) MineL1Block() {
	var blockID string
	s.Nil(s.RPCClient.L1.CallContext(context.Background(), &blockID, "evm_mine"))
	s.NotEmpty(blockID)
}
