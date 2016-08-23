package transaction

import (
	// Stdlib
	"io"
	"strings"

	// RPC
	"github.com/go-steem/rpc/apis/types"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (encoder *Encoder) Encode(v interface{}) error {
	if marshaller, ok := v.(TransactionMarshaller); ok {
		return marshaller.MarshalTransaction(encoder.w)
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
	case uint:
		return encoder.encodeUint(uint64(v))
	case uint8:
		return encoder.encodeUint(uint64(v))
	case uint16:
		return encoder.encodeUint(uint64(v))
	case uint32:
		return encoder.encodeUint(uint64(v))
	case uint64:
		return encoder.encodeUint(uint64(v))
	case *types.Int:
		return encoder.encodePtrTypesInt(v)

	case string:
		return encoder.encodeString(v)

	default:
		return errors.Errorf("encoder: unsupported type encountered")
	}
}

func (encoder *Encoder) encodeString(v string) error {
	if err := encoder.encodeInt(int64(len(v))); err != nil {
		return errors.Wrapf(err, "encoder: failed to write string: %v", v)
	}

	_, err := io.Copy(encoder.w, strings.NewReader(v))
	return errors.Wrapf(err, "encoder: failed to write string: %v", v)
}

func (encoder *Encoder) encodePtrTypesInt(v *types.Int) error {

}
