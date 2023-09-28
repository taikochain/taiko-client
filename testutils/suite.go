package testutils

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"net/url"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/suite"
	"github.com/taikoxyz/taiko-client/pkg/jwt"
	"github.com/taikoxyz/taiko-client/pkg/rpc"
	capacity "github.com/taikoxyz/taiko-client/prover/capacity_manager"
	"github.com/taikoxyz/taiko-client/prover/server"
)

type ClientTestSuite struct {
	suite.Suite
	testnetL1SnapshotID string
	RpcClient           *rpc.Client
	TestAddrPrivKey     *ecdsa.PrivateKey
	TestAddr            common.Address
	ProverEndpoints     []*url.URL
	proverServer        *server.ProverServer
}

func (s *ClientTestSuite) SetupTest() {
	// Default logger
	log.Root().SetHandler(
		log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stdout, log.TerminalFormat(true))),
	)

	if os.Getenv("LOG_LEVEL") != "" {
		level, err := log.LvlFromString(os.Getenv("LOG_LEVEL"))
		if err != nil {
			log.Crit("Invalid log level", "level", os.Getenv("LOG_LEVEL"))
		}
		log.Root().SetHandler(
			log.LvlFilterHandler(level, log.StreamHandler(os.Stdout, log.TerminalFormat(true))),
		)
	}

	testAddrPrivKey, err := crypto.ToECDSA(
		common.Hex2Bytes("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"),
	)
	s.Nil(err)

	s.TestAddrPrivKey = testAddrPrivKey
	s.TestAddr = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

	jwtSecret, err := jwt.ParseSecretFromFile(os.Getenv("JWT_SECRET"))
	s.Nil(err)
	s.NotEmpty(jwtSecret)

	rpcCli, err := rpc.NewClient(context.Background(), &rpc.ClientConfig{
		L1Endpoint:        os.Getenv("L1_NODE_WS_ENDPOINT"),
		L2Endpoint:        os.Getenv("L2_EXECUTION_ENGINE_WS_ENDPOINT"),
		TaikoL1Address:    common.HexToAddress(os.Getenv("TAIKO_L1_ADDRESS")),
		TaikoL2Address:    common.HexToAddress(os.Getenv("TAIKO_L2_ADDRESS")),
		TaikoTokenAddress: common.HexToAddress(os.Getenv("TAIKO_TOKEN_ADDRESS")),
		L2EngineEndpoint:  os.Getenv("L2_EXECUTION_ENGINE_AUTH_ENDPOINT"),
		JwtSecret:         string(jwtSecret),
		RetryInterval:     backoff.DefaultMaxInterval,
	})
	s.Nil(err)
	s.RpcClient = rpcCli

	l1ProverPrivKey, err := crypto.ToECDSA(common.Hex2Bytes(os.Getenv("L1_PROVER_PRIVATE_KEY")))
	s.Nil(err)

	s.ProverEndpoints = []*url.URL{LocalRandomProverEndpoint()}
	s.proverServer = NewTestProverServer(s, l1ProverPrivKey, capacity.New(1024, 100*time.Second), s.ProverEndpoints[0])

	tokenBalance, err := rpcCli.TaikoL1.GetTaikoTokenBalance(nil, crypto.PubkeyToAddress(l1ProverPrivKey.PublicKey))
	s.Nil(err)

	if tokenBalance.Cmp(common.Big0) == 0 {
		opts, err := bind.NewKeyedTransactorWithChainID(l1ProverPrivKey, rpcCli.L1ChainID)
		s.Nil(err)

		premintAmount, ok := new(big.Int).SetString(os.Getenv("PREMINT_TOKEN_AMOUNT"), 10)
		s.True(ok)
		s.True(premintAmount.Cmp(common.Big0) > 0)

		_, err = rpcCli.TaikoToken.Approve(opts, common.HexToAddress(os.Getenv("TAIKO_L1_ADDRESS")), premintAmount)
		s.Nil(err)

		tx, err := rpcCli.TaikoL1.DepositTaikoToken(opts, premintAmount)
		s.Nil(err)
		_, err = rpc.WaitReceipt(context.Background(), rpcCli.L1, tx)
		s.Nil(err)
	}

	s.Nil(rpcCli.L1RawRPC.CallContext(context.Background(), &s.testnetL1SnapshotID, "evm_snapshot"))
	s.NotEmpty(s.testnetL1SnapshotID)
}

func (s *ClientTestSuite) TearDownTest() {
	var revertRes bool
	s.Nil(s.RpcClient.L1RawRPC.CallContext(context.Background(), &revertRes, "evm_revert", s.testnetL1SnapshotID))
	s.True(revertRes)

	s.Nil(rpc.SetHead(context.Background(), s.RpcClient.L2RawRPC, common.Big0))
	s.Nil(s.proverServer.Shutdown(context.Background()))
}

func (s *ClientTestSuite) SetL1Automine(automine bool) {
	s.Nil(s.RpcClient.L1RawRPC.CallContext(context.Background(), nil, "evm_setAutomine", automine))
}
