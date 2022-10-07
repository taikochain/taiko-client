package producer

import (
	"errors"
	"math/big"
	"net/http"

	"github.com/taikochain/taiko-client/core/types"
	"github.com/taikochain/taiko-client/log"
)

var (
	errRpcdUnhealthy = errors.New("ZKEVM RPCD endpoint is unhealthy")
)

type ZkevmRpcdProducer struct {
	RpcdEndpoint string
}

func NewZkevmRpcdProducer(rpcdEndpoint string) (*ZkevmRpcdProducer, error) {
	resp, err := http.Get(rpcdEndpoint + "/health")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errRpcdUnhealthy
	}

	return &ZkevmRpcdProducer{RpcdEndpoint: rpcdEndpoint}, nil
}

// RequestProof implements the ProofProducer interface.
func (d *ZkevmRpcdProducer) RequestProof(
	opts *ProofRequestOptions,
	blockID *big.Int,
	header *types.Header,
	resultCh chan *ProofWithHeader,
) error {
	log.Info(
		"Request proof from ZKEVM RPCD service",
		"blockID", blockID,
		"height", header.Number,
		"hash", header.Hash(),
	)

	// TODO: call zkevm RPCD to get a proof.
	go func() {
		resultCh <- &ProofWithHeader{
			BlockID: blockID,
			Header:  header,
			ZkProof: []byte{0x00},
		}
	}()
	return nil
}
