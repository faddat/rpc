package networkbroadcast

import (
	// RPC
	"github.com/go-steem/rpc/apis/database"
)

var operations = map[string]int{
	database.OpTypeVote: 0,
}
