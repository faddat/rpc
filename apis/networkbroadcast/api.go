package networkbroadcast

import (
	// RPC
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
