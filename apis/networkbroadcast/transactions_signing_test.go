package networkbroadcast

import (
	// Stdlib
	"encoding/hex"
	"testing"

	// RPC
	"github.com/go-steem/rpc/encoding/wif"
)

var wifs = []string{
	"5JLw5dgQAx6rhZEgNN5C2ds1V47RweGshynFSWFbaMohsYsBvE8",
}

var privateKeys = make([][]byte, 0, len(wifs))

func init() {
	for _, v := range wifs {
		privKey, err := wif.Decode(v)
		if err != nil {
			panic(err)
		}
		privateKeys = append(privateKeys, privKey)
	}
}

func TestTransaction_Digest(t *testing.T) {
	expected := "ccbcb7d64444356654febe83b8010ca50d99edd0389d273b63746ecaf21adb92"

	digest, err := Digest(&tx, SteemChain)
	if err != nil {
		t.Error(err)
	}

	got := hex.EncodeToString(digest)
	if got != expected {
		t.Errorf("got %v, expected %v", got, expected)
	}
}

func TestTransaction_Sign(t *testing.T) {
	expected := "207a373c828b872d52a6946d31a22b4530e83a20bf89b4b71e29bbdffb4877c584232297ba3f22929506bf9706b2e67db9f99e517580cfcabeae492e2472d0a0dd"

	sigs, err := Sign(&tx, SteemChain, privateKeys)
	if err != nil {
		t.Error(err)
	}

	got := hex.EncodeToString(sigs[0])

	if got != expected {
		t.Errorf("\n\nexpected:\n%v\n\ngot:\n%v\n\n", expected, got)
	}
}
