package selector

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
	"github.com/taikoxyz/taiko-client/testutils"
)

type ProverSelectorTestSuite struct {
	testutils.ClientTestSuite
	s             *ETHFeeEOASelector
	proverAddress common.Address
}

func (s *ProverSelectorTestSuite) SetupTest() {
	s.ClientTestSuite.SetupTest()

	l1ProverPrivKey, err := crypto.ToECDSA(common.Hex2Bytes(os.Getenv("L1_PROVER_PRIVATE_KEY")))
	s.Nil(err)
	s.proverAddress = crypto.PubkeyToAddress(l1ProverPrivKey.PublicKey)

	protocolConfigs, err := s.RpcClient.TaikoL1.GetConfig(nil)
	s.Nil(err)

	s.s, err = NewETHFeeEOASelector(
		&protocolConfigs,
		s.RpcClient,
		common.HexToAddress(os.Getenv("TAIKO_L1_ADDRESS")),
		common.Big256,
		common.Big2,
		[]*url.URL{s.ProverEndpoints[0]},
		32,
		1*time.Minute,
		1*time.Minute,
	)
	s.Nil(err)
}

func TestProverSelectorTestSuite(t *testing.T) {
	suite.Run(t, new(ProverSelectorTestSuite))
}
