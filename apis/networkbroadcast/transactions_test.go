package networkbroadcast

import (
	// Stdlib
	"testing"
)

func TestTransaction_Serialize(t *testing.T) {
	tx := Transaction{
		RefBlockNum:    36029,
		RefBlockPrefix: 1164960351,
		Expiration:     &types.Time{time.},
	}

	expectedHex := "bd8c5fe26f45f179a8570100057865726f63057865726f6306706973746f6e1027"
}
