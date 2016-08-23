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

func (t *Transaction) Serialize() ([]byte, error) {
	// Write buffer.
	var b bytes.Buffer

	// Get ready for varints.
	varintBuffer := make([]byte, binary.MaxVarintLen64)

	putVarint := func(x int64) error {
		n := binary.PutVarint(varintBuffer, x)
		data := varintBuffer[:n]
		_, err := b.Write(data)
		return err
	}

	// Write ref_block_num.
	if err := binary.Write(&b, binary.LittleEndian, t.RefBlockNum); err != nil {
		return nil, errors.Wrapf(err, "networkbroadcast: failed to encode RefBlockNum: %v", t.RefBlockNum)
	}

	// Write ref_block_prefix.
	if err := binary.Write(&b, binary.LittleEndian, t.RefBlockPrefix); err != nil {
		return nil, errors.Wrapf(err, "networkbroadcast: failed to encode RefBlockPrefix: %v", t.RefBlockPrefix)
	}

	// Write expiration.
	timestamp := uint32(t.Expiration.Unix())
	if err := binary.Write(&b, binary.LittleEndian, timestamp); err != nil {
		return nil, errors.Wrapf(err, "networkbroadcast: failed to encode Expiration: %v (Unix)", timestamp)
	}

	// Write the number of operations.
	if err := putVarint(int64(len(t.Operations))); err != nil {
		return nil, errors.Wrap(err, "networkbroadcast: failed to encode Operations length")
	}

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
