package metrics

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/metrics/prometheus"
	"github.com/taikochain/taiko-client/cmd/flags"
	"github.com/urfave/cli/v2"
)

// Metrics
var (
	// Drvier
	DriverL1HeadHeightGuage    = metrics.NewRegisteredGauge("driver/l1Head/height", nil)
	DriverL1CurrentHeightGuage = metrics.NewRegisteredGauge("driver/l1Current/height", nil)
	DriverL2HeadIDGuage        = metrics.NewRegisteredGauge("driver/l2Head/id", nil)
	DriverL2FinalizedIDGuage   = metrics.NewRegisteredGauge("driver/l2Finalized/id", nil)

	// Proposer
	ProposerProposeEpochCounter    = metrics.NewRegisteredCounter("proposer/epoch", nil)
	ProposerProposedTxListsCounter = metrics.NewRegisteredCounter("proposer/proposed/txLists", nil)
	ProposerProposedTxsCounter     = metrics.NewRegisteredCounter("proposer/proposed/txs", nil)
	ProposerInvalidTxsCounter      = metrics.NewRegisteredCounter("proposer/invalid/txs", nil)

	// Prover
	ProverLatestFinalizedIDGauge      = metrics.NewRegisteredGauge("prover/lastFinalized/id", nil)
	ProverQueuedProofCounter          = metrics.NewRegisteredCounter("prover/proof/all/queued", nil)
	ProverQueuedValidProofCounter     = metrics.NewRegisteredCounter("prover/proof/valid/queued", nil)
	ProverQueuedInvalidProofCounter   = metrics.NewRegisteredCounter("prover/proof/invalid/queued", nil)
	ProverReceivedProofCounter        = metrics.NewRegisteredCounter("prover/proof/all/received", nil)
	ProverReceivedValidProofCounter   = metrics.NewRegisteredCounter("prover/proof/valid/received", nil)
	ProverReceivedInvalidProofCounter = metrics.NewRegisteredCounter("prover/proof/invalid/received", nil)
	ProverSentProofCounter            = metrics.NewRegisteredCounter("prover/proof/all/sent", nil)
	ProverSentValidProofCounter       = metrics.NewRegisteredCounter("prover/proof/valid/sent", nil)
	ProverSentInvalidProofCounter     = metrics.NewRegisteredCounter("prover/proof/invalid/sent", nil)
	ProverReceivedProposedBlockGauge  = metrics.NewRegisteredGauge("prover/proposed/received", nil)
)

// Serve starts the metrics server on the given address, will be close when the given
// context is canceled.
func Serve(ctx context.Context, c *cli.Context) error {
	if !c.Bool(flags.MetricsEnabled.Name) {
		return nil
	}

	address := net.JoinHostPort(
		c.String(flags.MetricsAddr.Name),
		strconv.Itoa(c.Int(flags.MetricsPort.Name)),
	)

	server := &http.Server{
		Addr:    address,
		Handler: prometheus.Handler(metrics.DefaultRegistry),
	}

	go func() {
		<-ctx.Done()
		if err := server.Close(); err != nil {
			log.Error("Failed to close metrics server", "error", err)
		}
	}()

	log.Info("Starting metrics server", "address", address)

	return server.ListenAndServe()
}
