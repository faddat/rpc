package networkbroadcast

import (
	// Stdlib
	"encoding/hex"
	"testing"
	"time"

	// RPC
	"github.com/go-steem/rpc/apis/database"
	"github.com/go-steem/rpc/encoding/wif"
	"github.com/go-steem/rpc/types"
)

func TestTransaction_Sign(t *testing.T) {
	var (
		WIF            = "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
		refBlockNum    = 34294
		refBlockPrefix = 3707022213
		expiration     = time.Date(2016, 4, 6, 8, 29, 27, 0, time.UTC)
	)

	tx := database.Transaction{
		RefBlockNum:    types.UInt16(refBlockNum),
		RefBlockPrefix: types.UInt32(refBlockPrefix),
		Expiration:     &types.Time{&expiration},
		Operations: []*database.Operation{
			{
				database.OpTypeVote,
				&database.VoteOperation{
					Voter:    "foobara",
					Author:   "foobarc",
					Permlink: "foobard",
					Weight:   1000,
				},
			},
		},
	}

	privKey, err := wif.Decode(WIF)
	if err != nil {
		t.Error(err)
	}

	privKeys := [][]byte{
		[]byte(privKey),
	}

	sigs, err := Sign(&tx, SteemChain, privKeys)
	if err != nil {
		t.Error(err)
	}

	sig := hex.EncodeToString(sigs[0])
	t.Log(sig)
}
