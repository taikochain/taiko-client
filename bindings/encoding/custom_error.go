package encoding

import (
	"context"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BlockHashContractCallerAndTransactionReader represents a contract caller and transaction reader.
type BlockHashContractCallerAndTransactionReader interface {
	bind.BlockHashContractCaller
	ethereum.TransactionReader
}

// TryParsingCustomErrorFromReceipt tries to parse the custom error from the given receipt.
func TryParsingCustomErrorFromReceipt(
	ctx context.Context,
	rpc BlockHashContractCallerAndTransactionReader,
	from common.Address,
	receipt *types.Receipt,
) error {
	// Fetch the raw transaction.
	tx, _, err := rpc.TransactionByHash(ctx, receipt.TxHash)
	if err != nil {
		return err
	}

	// Call the contract at the block hash.
	_, err = rpc.CallContractAtHash(ctx, ethereum.CallMsg{
		From:          from,
		To:            tx.To(),
		Gas:           tx.Gas(),
		GasPrice:      tx.GasPrice(),
		GasFeeCap:     tx.GasFeeCap(),
		GasTipCap:     tx.GasTipCap(),
		Value:         tx.Value(),
		Data:          tx.Data(),
		AccessList:    tx.AccessList(),
		BlobGasFeeCap: tx.BlobGasFeeCap(),
		BlobHashes:    tx.BlobHashes(),
	}, receipt.BlockHash)

	return TryParsingCustomError(err)
}

// TryParsingCustomError tries to checks whether the given error is one of the
// custom errors defined the protocol ABIs, if so, it will return
// the matched custom error, otherwise, it simply returns the original error.
func TryParsingCustomError(originalError error) error {
	if originalError == nil {
		return nil
	}

	errData := getErrorData(originalError)

	// if errData is unparsable and returns 0x, we should not match any errors.
	if errData == "0x" {
		return originalError
	}

	for _, customErrors := range customErrorMaps {
		for _, customError := range customErrors {
			if strings.HasPrefix(customError.ID.Hex(), errData) {
				return errors.New(customError.Name)
			}
		}
	}

	return originalError
}

// getErrorData tries to parse the actual custom error data from the given error.
func getErrorData(err error) string {
	// Geth node custom errors, the actual struct of this error is go-ethereum's <rpc.jsonError Value>.
	gethJSONError, ok := err.(interface{ ErrorData() interface{} }) // nolint: errorlint
	if ok {
		if errData, ok := gethJSONError.ErrorData().(string); ok {
			return errData
		}
	}

	// Hardhat node custom errors, example:
	// "VM Exception while processing transaction: reverted with an unrecognized custom error (return data: 0xb6d363fd)"
	if strings.Contains(err.Error(), "reverted with an unrecognized custom error") {
		return err.Error()[len(err.Error())-11 : len(err.Error())-1]
	}

	return err.Error()
}
