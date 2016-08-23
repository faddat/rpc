package transaction

type TransactionMarshaller interface {
	MarshalTransaction(io.Writer) error
}
