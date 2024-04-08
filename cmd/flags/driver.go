package flags

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/urfave/cli/v2"
)

// Optional flags used by driver.
var (
	SyncModeFlag = &cli.StringFlag{
		Name:     "syncmode",
		Usage:    `Blockchain sync mode ("snap" or "full")`,
		Value:    "snap",
		Category: driverCategory,
		Action: func(_ *cli.Context, s string) error {
			if s != downloader.SnapSync.String() &&
				s != downloader.FullSync.String() {
				return fmt.Errorf("invalid sync mode: %s", s)
			}
			return nil
		},
	}
	P2PSyncVerifiedBlocks = &cli.BoolFlag{
		Name: "p2p.syncVerifiedBlocks",
		Usage: "Try P2P syncing verified blocks between L2 execution engines, " +
			"will be helpful to bring a new node online quickly",
		Value:    false,
		Category: driverCategory,
	}
	P2PSyncTimeout = &cli.DurationFlag{
		Name: "p2p.syncTimeout",
		Usage: "P2P syncing timeout, if no sync progress is made within this time span, " +
			"driver will stop the P2P sync and insert all remaining L2 blocks one by one",
		Value:    1 * time.Hour,
		Category: driverCategory,
	}
	CheckPointSyncURL = &cli.StringFlag{
		Name:     "p2p.checkPointSyncUrl",
		Usage:    "HTTP RPC endpoint of another synced L2 execution engine node",
		Category: driverCategory,
	}
	// syncer specific flag
	MaxExponent = &cli.Uint64Flag{
		Name: "syncer.maxExponent",
		Usage: "Maximum exponent of retrieving L1 blocks when there is a mismatch between protocol and L2 EE," +
			"0 means that it is reset to the genesis height",
		Value:    0,
		Category: driverCategory,
	}
)

// DriverFlags All driver flags.
var DriverFlags = MergeFlags(CommonFlags, []cli.Flag{
	L1BeaconEndpoint,
	L2WSEndpoint,
	L2AuthEndpoint,
	JWTSecret,
	SyncModeFlag,
	P2PSyncVerifiedBlocks,
	P2PSyncTimeout,
	CheckPointSyncURL,
	MaxExponent,
})
