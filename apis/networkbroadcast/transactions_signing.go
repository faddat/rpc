// +build !nosigning

package networkbroadcast

import (
	// Stdlib
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"unsafe"

	// RPC
	"github.com/go-steem/rpc/apis/database"

	// Vendor
	"github.com/pkg/errors"
)

// #cgo LDFLAGS: -lsecp256k1
// #include <stdlib.h>
// #include "signing.h"
import "C"

func Sign(tx *database.Transaction, chain *Chain, privKeys [][]byte) ([][]byte, error) {
	var messageBuffer bytes.Buffer

	// Write the chain ID.
	rawChainID, err := hex.DecodeString(chain.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode chain ID: %v", chain.ID)
	}

	if _, err := messageBuffer.Write(rawChainID); err != nil {
		return nil, errors.Wrap(err, "failed to write chain ID")
	}

	// Write the serialized transaction.
	rawTx, err := Serialize(tx)
	if err != nil {
		return nil, err
	}

	if _, err := messageBuffer.Write(rawTx); err != nil {
		return nil, errors.Wrap(err, "failed to write serialized transaction")
	}

	// Compute the digest.
	digest := sha256.Sum256(messageBuffer.Bytes())

	// Sign.
	cDigest := C.CBytes(digest[:])
	defer C.free(cDigest)

	cKeys := make([]unsafe.Pointer, 0, len(privKeys))
	for _, key := range privKeys {
		cKeys = append(cKeys, C.CBytes(key))
	}
	defer func() {
		for _, cKey := range cKeys {
			C.free(cKey)
		}
	}()

	sigs := make([][]byte, 0, len(privKeys))
	for _, cKey := range cKeys {
		signature := make([]byte, 64)
		var recid C.int

		i := int(C.sign_transaction(
			(*C.uchar)(cDigest), (*C.uchar)(cKey), (*C.uchar)(&signature[0]), &recid))
		if i == 0 {
			return nil, errors.New("sign_transaction returned a non-zero exit status")
		}

		sig := make([]byte, 65)
		sig[0] = byte(recid)
		copy(sig[1:], signature)

		sigs = append(sigs, sig)
	}

	return sigs, nil
}
