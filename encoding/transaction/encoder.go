package transaction

import (
	// Stdlib
	"encoding/binary"
	"io"
	"strings"

	// RPC
	"github.com/go-steem/rpc/apis/types"

	// Vendor
	"github.com/pkg/errors"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (encoder *Encoder) Encode(v interface{}) error {
	if marshaller, ok := v.(TransactionMarshaller); ok {
		return marshaller.MarshalTransaction(encoder)
	}

	switch v := v.(type) {
	case int:
		return encoder.encodeInt(int64(v))
	case int8:
		return encoder.encodeInt(int64(v))
	case int16:
		return encoder.encodeInt(int64(v))
	case int32:
		return encoder.encodeInt(int64(v))
	case int64:
		return encoder.encodeInt(int64(v))

	case types.Int8:
		return encoder.encodeInt(int64(v))
	case types.Int16:
		return encoder.encodeInt(int64(v))
	case types.Int32:
		return encoder.encodeInt(int64(v))
	case types.Int64:
		return encoder.encodeInt(int64(v))

	case uint:
		return encoder.encodeUInt(uint64(v))
	case uint8:
		return encoder.encodeUInt(uint64(v))
	case uint16:
		return encoder.encodeUInt(uint64(v))
	case uint32:
		return encoder.encodeUInt(uint64(v))
	case uint64:
		return encoder.encodeUInt(uint64(v))

	case types.UInt:
		return encoder.encodeUInt(uint64(v))
	case types.UInt8:
		return encoder.encodeUInt(uint64(v))
	case types.UInt16:
		return encoder.encodeUInt(uint64(v))
	case types.UInt32:
		return encoder.encodeUInt(uint64(v))
	case types.UInt64:
		return encoder.encodeUInt(uint64(v))

	case string:
		return encoder.encodeString(v)

	default:
		return errors.Errorf("encoder: unsupported type encountered")
	}
}

func (encoder *Encoder) encodeInt(v int64) error {
	if v >= 0 {
		return encoder.encodeUInt(uint64(v))
	}

	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(b, v)
	return encoder.writeBytes(b[:n])
}

func (encoder *Encoder) encodeUInt(v uint64) error {
	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(b, v)
	return encoder.writeBytes(b[:n])
}

func (encoder *Encoder) encodeString(v string) error {
	if err := encoder.encodeUInt(uint64(len(v))); err != nil {
		return errors.Wrapf(err, "encoder: failed to write string: %v", v)
	}

	return encoder.writeString(v)
}

func (encoder *Encoder) writeBytes(bs []byte) error {
	if _, err := encoder.w.Write(bs); err != nil {
		return errors.Wrapf(err, "encoder: failed to write bytes: %v", bs)
	}
	return nil
}

func (encoder *Encoder) writeString(s string) error {
	if _, err := io.Copy(encoder.w, strings.NewReader(s)); err != nil {
		return errors.Wrapf(err, "encoder: failed to write string: %v", s)
	}
	return nil
}
