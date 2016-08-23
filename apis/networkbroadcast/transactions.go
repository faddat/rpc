package networkbroadcast

import (
	// Stdlib
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	// RPC
	"github.com/go-steem/rpc/apis/database"
	"github.com/go-steem/rpc/encoding/transaction"

	// Vendor
	"github.com/pkg/errors"
)

func Serialize(tx *database.Transaction) ([]byte, error) {
	// Prepare an encoder.
	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	// Write ref_block_num.
	if err := encoder.Encode(tx.RefBlockNum); err != nil {
		return nil, errors.Wrapf(err,
			"networkbroadcast: failed to encode RefBlockNum: %v", tx.RefBlockNum)
	}

	// Write ref_block_prefix.
	if err := encoder.Encode(tx.RefBlockPrefix); err != nil {
		return nil, errors.Wrapf(err,
			"networkbroadcast: failed to encode RefBlockPrefix: %v", tx.RefBlockPrefix)
	}

	// Write expiration.
	if err := encoder.Encode(tx.Expiration); err != nil {
		return nil, errors.Wrapf(err,
			"networkbroadcast: failed to encode Expiration: %v", tx.Expiration)
	}

	// Write the number of operations.
	if err := encoder.EncodeUVarint(uint64(len(tx.Operations))); err != nil {
		return nil, errors.Wrap(err,
			"networkbroadcast: failed to encode Operations length")
	}

	// Write the operations, one by one.
	for _, op := range tx.Operations {
		if err := encoder.Encode(op.Body); err != nil {
			return nil, errors.Wrap(err, "networkbroadcast: failed to encode an operation")
		}
	}

	// Return the result.
	return b.Bytes(), nil
}

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

func RefBlockNum(blockNumber uint32) uint16 {
	return uint16(blockNumber)
}

func RefBlockPrefix(blockID string) (uint32, error) {
	// Block ID is hex-encoded.
	rawBlockID, err := hex.DecodeString(blockID)
	if err != nil {
		return 0, errors.Wrapf(err, "networkbroadcast: failed to decode block ID: %v", blockID)
	}

	// Raw prefix = raw block ID [4:8].
	// Make sure we don't trigger a slice bounds out of range panic.
	if len(rawBlockID) < 8 {
		return 0, errors.Errorf("networkbroadcast: invalid block ID: %v", blockID)
	}
	rawPrefix := rawBlockID[4:8]

	// Decode the prefix.
	var prefix uint32
	if err := binary.Read(bytes.NewReader(rawPrefix), binary.LittleEndian, &prefix); err != nil {
		return 0, errors.Wrapf(err, "networkbroadcast: failed to read block prefix: %v", rawPrefix)
	}

	// Done, return the prefix.
	return prefix, nil
}
