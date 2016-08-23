package transaction

import "io"

type TransactionMarshaller interface {
	MarshalTransaction(io.Writer) error
}
