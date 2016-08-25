package wif

import (
	// Stdlib
	"encoding/hex"
	"testing"
)

func TestDecode(t *testing.T) {
	privKey, err := Decode(wif)
	if err != nil {
		t.Error(err)
	}

	expected := privKeyHex
	got := hex.EncodeToString(privKey)

	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
