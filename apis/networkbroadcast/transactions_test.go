package networkbroadcast

import (
	// Stdlib
	"encoding/hex"
	"testing"
	"time"

	// RPC
	"github.com/go-steem/rpc/apis/database"
	"github.com/go-steem/rpc/types"
)

var tx database.Transaction

func init() {
	expiration := time.Date(2016, 8, 8, 12, 24, 17, 0, time.UTC)

	tx = database.Transaction{
		RefBlockNum:    36029,
		RefBlockPrefix: 1164960351,
		Expiration:     &types.Time{&expiration},
		Operations: []*database.Operation{
			{
				database.OpTypeVote,
				&database.VoteOperation{
					Voter:    "xeroc",
					Author:   "xeroc",
					Permlink: "piston",
					Weight:   10000,
				},
			},
		},
	}
}

func TestTransaction_Serialize(t *testing.T) {
	expectedHex := "bd8c5fe26f45f179a8570100057865726f63057865726f6306706973746f6e1027"

	serialized, err := Serialize(&tx)
	if err != nil {
		t.Error(err)
	}

	serializedHex := hex.EncodeToString(serialized)

	if serializedHex != expectedHex {
		t.Errorf("expected %v, got %v", expectedHex, serializedHex)
	}
}
