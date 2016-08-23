package transaction

import "io"

type TransactionMarshaller interface {
	MarshalTransaction(*Encoder) error
}
