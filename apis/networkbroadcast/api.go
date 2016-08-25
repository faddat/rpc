package networkbroadcast

import (
	// Stdlib
	"encoding/json"

	// RPC
	"github.com/go-steem/rpc/apis/database"
	"github.com/go-steem/rpc/interfaces"
)

type API struct {
	caller interfaces.Caller
}

func NewAPI(caller interfaces.Caller) *API {
	return &API{caller}
}

func (api *API) call(method string, params, resp interface{}) error {
	return api.caller.Call("call", []interface{}{"network_broadcast_api", method, params}, resp)
}

func (api *API) BroadcastTransaction(tx *database.Transaction) error {
	return api.call("broadcast_transaction", tx, nil)
}

func (api *API) BroadcastTransactionSynchronousRaw(tx *database.Transaction) (*json.RawMessage, error) {
	var resp json.RawMessage
	if err := api.call("broadcast_transaction", tx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
