package networkbroadcast

import (
	// Stdlib
	"encoding/json"

	// RPC
	"github.com/go-steem/rpc/interfaces"
	"github.com/go-steem/rpc/types"
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

func (api *API) BroadcastTransaction(tx *types.Transaction) error {
	return api.call("broadcast_transaction", tx, nil)
}

func (api *API) BroadcastTransactionSynchronousRaw(tx *types.Transaction) (*json.RawMessage, error) {
	var resp json.RawMessage
	params := []interface{}{tx}
	if err := api.call("broadcast_transaction_synchronous", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
