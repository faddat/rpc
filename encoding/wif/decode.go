package wif

import (
	// Vendor
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// Decode can be used to turn WIF into a raw private key (32 bytes).
func Decode(wif string) ([]byte, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode WIF")
	}

	return w.PrivKey.Serialize(), nil
}
