package customjson

import (
	// Stdlib
	"encoding/json"
	"reflect"
	"strings"

	// Vendor
	"github.com/pkg/errors"
)

const (
	TypeFollow = "follow"
)

var dataObjects = map[string]interface{}{
	TypeFollow: &FollowOperation{},
}

// FC_REFLECT( steemit::chain::custom_json_operation,
//             (required_auths)
//             (required_posting_auths)
//             (id)
//             (json) )

// Operation represents custom_json operation data.
type Operation struct {
	RequiredAuths        []string `json:"required_auths"`
	RequiredPostingAuths []string `json:"required_posting_auths"`
	ID                   string   `json:"id"`
	JSON                 string   `json:"json"`
}

func (op *Operation) UnmarshalData() (interface{}, error) {
	// Get the corresponding data object template.
	template, ok := dataObjects[op.ID]
	if !ok {
		// In case there is no corresponding template, return nil.
		return nil, nil
	}

	// Clone the template.
	opData := reflect.New(reflect.Indirect(reflect.ValueOf(template)).Type()).Interface()

	// Unmarshal into the newly created data object instance.
	if err := json.NewDecoder(strings.NewReader(op.JSON)).Decode(opData); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal custom_json operation data")
	}

	return opData, nil
}
