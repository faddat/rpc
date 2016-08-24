package networkbroadcast

import (
	// Stdlib
	"encoding/hex"
	"testing"
)

func TestTransaction_Sign(t *testing.T) {
	privKeys := [][]byte{
		[]byte("5JLw5dgQAx6rhZEgNN5C2ds1V47RweGshynFSWFbaMohsYsBvE8"),
	}

	sigs, err := Sign(&tx, SteemChain, privKeys)
	if err != nil {
		t.Error(err)
	}

	sig := hex.EncodeToString(sigs[0])
	t.Log(sig)
}
