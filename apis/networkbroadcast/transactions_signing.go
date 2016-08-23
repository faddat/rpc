// +build !nosigning

package networkbroadcast

import (
	// Stdlib
	"bytes"
	"crypto/sha256"
	"encoding/hex"

	// RPC
	"github.com/go-steem/rpc/apis/database"

	// Vendor
	"github.com/pkg/errors"
)

func Sign(tx *database.Transaction, chain *Chain, privKeys [][]byte) ([]byte, error) {
	var messageBuffer bytes.Buffer

	// Write the chain ID.
	rawChainID, err := hex.DecodeString(chain.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode chain ID: %v", chain.ID)
	}

	if _, err := messageBuffer.Write(rawChainID); err != nil {
		return nil, errors.Wrap(err, "failed to write chain ID")
	}

	// Write the serialized transaction.
	rawTx, err := Serialize(tx)
	if err != nil {
		return nil, err
	}

	if _, err := messageBuffer.Write(rawTx); err != nil {
		return nil, errors.Wrap(err, "failed to write serialized transaction")
	}

	// Compute the digest.
	digest := sha256.Sum256(messageBuffer.Bytes())

	panic("Not Implemented")
	return nil, nil
}
