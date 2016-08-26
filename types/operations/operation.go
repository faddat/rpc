package operations

import (
	// Stdlib
	"encoding/json"
	"reflect"

	// RPC
	"github.com/go-steem/rpc/types/operations/customjson"

	// Vendor
	"github.com/pkg/errors"
)

// dataObjects keeps mapping operation type -> operation data object.
// This is used later on to unmarshal operation data based on the operation type.
var dataObjects = map[OpType]interface{}{
	TypeVote:              &VoteOperation{},
	TypeComment:           &CommentOperation{},
	TypeTransfer:          &TransferOperation{},
	TypeTransferToVesting: &TransferToVestingOperation{},
	TypeWithdrawVesting:   &WithdrawVestingOperation{},
	TypeLimitOrderCreate:  &LimitOrderCreateOperation{},
	TypeLimitOrderCancel:  &LimitOrderCancelOperation{},
	TypeFeedPublish:       &FeedPublishOperation{},
	TypeConvert:           &ConvertOperation{},
	TypeAccountCreate:     &AccountCreateOperation{},
	TypeAccountUpdate:     &AccountUpdateOperation{},
	// TypeWitnessUpdate:  &WitnessUpdateOperation{},
	TypeAccountWitnessVote:  &AccountWitnessVoteOperation{},
	TypeAccountWitnessProxy: &AccountWitnessProxyOperation{},
	TypePOW:                 &POWOperation{},
	// TypeCustom:           &CustomOperation{},
	TypeReportOverProduction: &ReportOverProductionOperation{},
	TypeDeleteComment:        &DeleteCommentOperation{},
	TypeCustomJSON:           &customjson.Operation{},
	TypeCommentOptions:       &CommentOptionsOperation{},
	// TypeSetWithdrawVestingRoute: &SetWithdrawVestingRouteOperation{},
	// TypeLimitOrderCreate2:       &LimitOrderCreate2Operation{},
	// TypeChallengeAuthority:      &ChallengeAuthorityOperation{},
	// TypeProveAuthority:          &ProveAuthorityOperation{},
	// TypeRequestAccountRecovery:  &RequestAccountRecoveryOperation{},
	// TypeRecoverAccount:          &RecoverAccountOperation{},
	// TypeChangeRecoveryAccount:   &ChangeRecoverAccountOperation{},
	// TypeEscrowTransfer:          &EscrowTransferOperation{},
	// TypeEscrowDispute:           &EscrowDisputeOperation{},
	// TypeEscrowRelease:           &EescrowReleaseOperation{},
	// TypePOW2:                    &POW2Operation{},
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

func (op *Operation) MarshalJSON() ([]byte, error) {
	tuple := []interface{}{
		op.Type,
		op.Data,
	}
	return json.Marshal(typle)
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
