package networkbroadcast

import (
	// Stdlib
	"bytes"
	"encoding/binary"
	"encoding/hex"

	// RPC
	"github.com/go-steem/rpc/apis/types"

	// Vendor
	"github.com/pkg/errors"
)

type Transaction struct {
	RefBlockNum    uint16
	RefBlockPrefix uint32
	Operations     []interface{}
	Expiration     *types.Time
}

func (tx *Transaction) Serialize() ([]byte, error) {
	// Prepare an encoder.
	var b bytes.Buffer
	encoder := transactions.NewEncoder(&w)

	// Write ref_block_num.
	if err := encoder.Encode(tx.RefBlockNum); err != nil {
		return nil, errors.Wrapf(err, "networkbroadcast: failed to encode RefBlockNum: %v", tx.RefBlockNum)
	}

	// Write ref_block_prefix.
	if err := encoder.Encode(tx.RefBlockPrefix); err != nil {
		return nil, errors.Wrapf(err, "networkbroadcast: failed to encode RefBlockPrefix: %v", tx.RefBlockPrefix)
	}

	// Write expiration.
	timestamp := uint32(tx.Expiration.Unix())
	if err := encoder.Encode(timestamp); err != nil {
		return nil, errors.Wrapf(err, "networkbroadcast: failed to encode Expiration: %v (Unix)", timestamp)
	}

	// Write the number of operations.
	if err := encoder.EncodeUVarint(uint64(len(tx.Operation))); err != nil {
		return nil, errors.Wrap(err, "networkbroadcast: failed to encode Operations length")
	}

	// Write the operations, one by one.
	for _, op := range tx.Operations {
		if err := encoder.Encode(op); err != nil {
			return nil, errors.Wrap(err, "networkbroadcast: failed to encode an operation")
		}
	}

	// Return the result.
	return b.Bytes(), nil
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
