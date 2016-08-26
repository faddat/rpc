package types

import (
	// RPC
	"github.com/go-steem/rpc/encoding/transaction"
)

// Transaction represents a blockchain transaction.
type Transaction struct {
	RefBlockNum    UInt16
	RefBlockPrefix UInt32
	Expiration     *Time
	Operations     []*Operation
	Signatures     []string
}

// MarshalTransaction implements transaction.Marshaller interface.
func (tx *Transaction) MarshalTransaction(encoder *transaction.Encoder) error {
	if len(tx.Operations) == 0 {
		return errors.New("no operation specified")
	}

	enc := transaction.NewRollingEncoder(encoder)

	enc.Encode(tx.RefBlockNum)
	enc.Encode(tx.RefBlockPrefix)
	enc.Encode(tx.Expiration)

	enc.EncodeUVarint(uint64(len(tx.Operations)))
	for _, op := range tx.Operations {
		enc.Encode(op.Data)
	}

	return enc.Err()
}
