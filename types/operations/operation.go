package operations

import (
	// Stdlib
	"encoding/json"
	"reflect"

	// Vendor
	"github.com/pkg/errors"
)

// dataObjects keeps mapping operation type -> operation data object.
// This is used later on to unmarshal operation data based on the operation type.
var dataObjects = map[OpType]interface{}{
	TypeConvert:             &ConvertOperation{},
	TypeFeedPublish:         &FeedPublishOperation{},
	TypePOW:                 &PowOperation{},
	TypeCustomJSON:          &CustomJSONOperation{},
	TypeAccountCreate:       &AccountCreateOperation{},
	TypeAccountUpdate:       &AccountUpdateOperation{},
	TypeTransfer:            &TransferOperation{},
	TypeTransferToVesting:   &TransferToVestingOperation{},
	TypeWithdrawVesting:     &WithdrawVestingOperation{},
	TypeAccountWitnessVote:  &AccountWitnessVoteOperation{},
	TypeAccountWitnessProxy: &AccountWitnessProxyOperation{},
	TypeComment:             &CommentOperation{},
	TypeVote:                &VoteOperation{},
	TypeLimitOrderCreate:    &LimitOrderCreateOperation{},
	TypeLimitOrderCancel:    &LimitOrderCancelOperation{},
	TypeDeleteComment:       &DeleteCommentOperation{},
	TypeCommentOptions:      &CommentOptionsOperation{},
}

// Operation represents an operation stored in a transaction.
type Operation struct {
	// Type contains the operation type as present in the operation object, element [0].
	Type OpType

	// Data contains the operation data as present in the operation object, element [1].
	//
	// When the operation type is known to this package, this field contains
	// the operation data object associated with the given operation type,
	// e.g. Type is TypeVote -> Data contains *VoteOperation.
	// Otherwise this field contains raw JSON (type *json.RawMessage).
	Data interface{}
}

func (op *Operation) UnmarshalJSON(data []byte) error {
	// The operation object is [opType, opBody].
	raw := make([]*json.RawMessage, 2)
	if err := json.Unmarshal(data, &raw); err != nil {
		return errors.Wrapf(err, "failed to unmarshal operation object: %v", string(data))
	}
	if len(raw) != 2 {
		return errors.Errorf("invalid operation object: %v", string(data))
	}

	// Unmarshal the type.
	var opType OpType
	if err := json.Unmarshal(*raw[0], &opType); err != nil {
		return errors.Wrapf(err, "failed to unmarshal Operation.Type: %v", string(*raw[0]))
	}

	// Unmarshal the data.
	var opData interface{}
	template, ok := dataObjects[opType]
	if ok {
		opData = reflect.New(reflect.Indirect(reflect.ValueOf(template)).Type()).Interface()
		if err := json.Unmarshal(*raw[1], opData); err != nil {
			return errors.Wrapf(err, "failed to unmarshal Operation.Data: %v", string(*raw[1]))
		}
	} else {
		opData = raw[1]
	}

	// Update fields.
	op.Type = opType
	op.Data = opData
	return nil
}
