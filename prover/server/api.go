package server

import (
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/labstack/echo/v4"

	"github.com/taikoxyz/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-client/pkg/rpc"
)

// @title Taiko Prover Server API
// @version 1.0
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://community.taiko.xyz/
// @contact.email info@taiko.xyz

// @license.name MIT
// @license.url https://github.com/taikoxyz/taiko-client/blob/main/LICENSE.md

// CreateAssignmentRequestBody represents a request body when handling assignment creation request.
type CreateAssignmentRequestBody struct {
	FeeToken   common.Address
	TierFees   []encoding.TierFee
	Expiry     uint64
	TxListHash common.Hash
}

// Status represents the current prover server status.
type Status struct {
	MinOptimisticTierFee uint64 `json:"minOptimisticTierFee"`
	MinSgxTierFee        uint64 `json:"minSgxTierFee"`
	MinSgxAndZkVMTierFee uint64 `json:"minSgxAndZkVMTierFee"`
	MaxExpiry            uint64 `json:"maxExpiry"`
	Prover               string `json:"prover"`
}

// GetStatus handles a query to the current prover server status.
//
//	@Summary		Get current prover server status
//	@ID			   	get-status
//	@Accept			json
//	@Produce		json
//	@Success		200	{object} Status
//	@Router			/status [get]
func (s *ProverServer) GetStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, &Status{
		MinOptimisticTierFee: s.minOptimisticTierFee.Uint64(),
		MinSgxTierFee:        s.minSgxTierFee.Uint64(),
		MinSgxAndZkVMTierFee: s.minSgxAndZkVMTierFee.Uint64(),
		MaxExpiry:            uint64(s.maxExpiry.Seconds()),
		Prover:               s.proverAddress.Hex(),
	})
}

// ProposeBlockResponse represents the JSON response which will be returned by
// the ProposeBlock request handler.
type ProposeBlockResponse struct {
	SignedPayload []byte         `json:"signedPayload"`
	Prover        common.Address `json:"prover"`
	MaxBlockID    uint64         `json:"maxBlockID"`
	MaxProposedIn uint64         `json:"maxProposedIn"`
}

// CreateAssignment handles a block proof assignment request, decides if this prover wants to
// handle this block, and if so, returns a signed payload the proposer
// can submit onchain.
//
//	@Summary		Try to accept a block proof assignment
//	@Param          body        body    CreateAssignmentRequestBody   true    "assignment request body"
//	@Accept			json
//	@Produce		json
//	@Success		200		{object} ProposeBlockResponse
//	@Failure		422		{string} string	"invalid txList hash"
//	@Failure		422		{string} string	"only receive ETH"
//	@Failure		422		{string} string	"insufficient prover balance"
//	@Failure		422		{string} string	"proof fee too low"
//	@Failure		422		{string} string "expiry too long"
//	@Failure		422		{string} string "prover does not have capacity"
//	@Router			/assignment [post]
func (s *ProverServer) CreateAssignment(c echo.Context) error {
	req := new(CreateAssignmentRequestBody)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	log.Info(
		"Proof assignment request body",
		"feeToken", req.FeeToken,
		"expiry", req.Expiry,
		"tierFees", req.TierFees,
		"txListHash", req.TxListHash,
		"currentUsedCapacity", len(s.proofSubmissionCh),
	)

	if req.TxListHash == (common.Hash{}) {
		log.Info("Invalid txList hash")
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "invalid txList hash")
	}

	if req.FeeToken != (common.Address{}) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "only receive ETH")
	}

	ok, err := rpc.CheckProverBalance(
		c.Request().Context(),
		s.rpc,
		s.proverAddress,
		s.assignmentHookAddress,
		s.livenessBond,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if !ok {
		log.Warn(
			"Insufficient prover balance, please get more tokens or wait for verification of the blocks you proved",
			"prover", s.proverAddress,
		)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "insufficient prover balance")
	}

	for _, tier := range req.TierFees {
		if tier.Tier == encoding.TierGuardianID {
			continue
		}

		var minTierFee *big.Int
		switch tier.Tier {
		case encoding.TierOptimisticID:
			minTierFee = s.minOptimisticTierFee
		case encoding.TierSgxID:
			minTierFee = s.minSgxTierFee
		case encoding.TierSgxAndZkVMID:
			minTierFee = s.minSgxAndZkVMTierFee
		default:
			log.Warn("Unknown tier", "tier", tier.Tier, "fee", tier.Fee, "proposerIP", c.RealIP())
		}

		if tier.Fee.Cmp(minTierFee) < 0 {
			log.Warn(
				"Proof fee too low",
				"tier", tier.Tier,
				"fee", tier.Fee,
				"minTierFee", minTierFee,
				"proposerIP", c.RealIP(),
			)
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "proof fee too low")
		}
	}

	if req.Expiry > uint64(time.Now().Add(s.maxExpiry).Unix()) {
		log.Warn(
			"Expiry too long",
			"requestExpiry", req.Expiry,
			"srvMaxExpiry", s.maxExpiry,
			"proposerIP", c.RealIP(),
		)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "expiry too long")
	}

	// Check if the prover has any capacity now.
	if s.proofSubmissionCh != nil && len(s.proofSubmissionCh) == cap(s.proofSubmissionCh) {
		log.Warn("Prover does not have capacity", "capacity", cap(s.proofSubmissionCh))
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "prover does not have capacity")
	}

	l1Head, err := s.rpc.L1.BlockNumber(c.Request().Context())
	if err != nil {
		log.Error("Failed to get L1 block head", "error", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	encoded, err := encoding.EncodeProverAssignmentPayload(
		s.protocolConfigs.ChainId,
		s.taikoL1Address,
		s.assignmentHookAddress,
		req.TxListHash,
		req.FeeToken,
		req.Expiry,
		l1Head+s.maxSlippage,
		s.maxProposedIn,
		req.TierFees,
	)
	if err != nil {
		log.Error("Failed to encode proverAssignment payload data", "error", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	signed, err := crypto.Sign(crypto.Keccak256Hash(encoded).Bytes(), s.proverPrivateKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &ProposeBlockResponse{
		SignedPayload: signed,
		Prover:        s.proverAddress,
		MaxBlockID:    l1Head + s.maxSlippage,
		MaxProposedIn: s.maxProposedIn,
	})
}
