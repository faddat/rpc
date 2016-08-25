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

func Digest(tx *database.Transaction, chain *Chain) ([]byte, error) {
	var msgBuffer bytes.Buffer

	// Write the chain ID.
	rawChainID, err := hex.DecodeString(chain.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode chain ID: %v", chain.ID)
	}

	if _, err := msgBuffer.Write(rawChainID); err != nil {
		return nil, errors.Wrap(err, "failed to write chain ID")
	}

	// Write the serialized transaction.
	rawTx, err := Serialize(tx)
	if err != nil {
		return nil, err
	}

	if _, err := msgBuffer.Write(rawTx); err != nil {
		return nil, errors.Wrap(err, "failed to write serialized transaction")
	}

	// Compute the digest.
	digest := sha256.Sum256(msgBuffer.Bytes())
	return digest[:], nil
}

func Sign(tx *database.Transaction, chain *Chain, privKeys [][]byte) ([][]byte, error) {
	digest, err := Digest(tx, chain)
	if err != nil {
		return nil, err
	}

	// Sign.
	cDigest := C.CBytes(digest)
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
		var (
			signature [64]byte
			recid     C.int
		)

		code := C.sign_transaction(
			(*C.uchar)(cDigest), (*C.uchar)(cKey), (*C.uchar)(&signature[0]), &recid)
		if code == 0 {
			return nil, errors.New("sign_transaction returned 0")
		}

		sig := make([]byte, 65)
		sig[0] = byte(recid)
		copy(sig[1:], signature[:])

		sigs = append(sigs, sig)
	}

	return sigs, nil
}

func Verify(tx *database.Transaction, chain *Chain, pubKeys [][]byte) (bool, error) {
	// Compute the digest, again.
	digest, err := Digest(tx)
	if err != nil {
		return false, err
	}

	cDigest := C.CBytes(digest)
	defer C.free(cDigest)

	// Make sure to free memory.
	cSigs := make([]unsafe.Pointer, 0, len(tx.Signatures))
	defer func() {
		for _, cSig := range cSigs {
			C.free(cSig)
		}
	}()

	// Collect verified public keys.
	pubKeysFound := make([][]byte, len(pubKeys))
	for i, signature := range tx.Signatures {
		sig, err := hex.DecodeString(signature)
		if err != nil {
			return false, errors.Wrap(err, "failed to decode signature hex")
		}

		recoverParameter := sig[0]
		sig := sig[1:]

		cSig := C.CBytes(sig)
		cSigs = append(cSigs, cSig)

		var publicKey [33]byte

		code := C.verify_recoverable_signature(
			(*C.char)(cDigest),
			(*C.char)(cSig),
			(C.int)(recoverParameter),
			(*C.char)(&publicKey[0]),
		)
		if code == 1 {
			pubKeysFound[i] = publicKey[:]
		}
	}

	for i := range pubKeys {
		if !bytes.Equal(pubKeysFound[i], pubKeys[i]) {
			return false, nil
		}
	}
	return true, nil
}

func verifySignature(pubKey []byte, message []byte, signature []byte) (bool, error) {
	code := C.verify_signature((*C.char)(cPubKey), (*C.char)(cMessage), (*C.char)(cSignature))
	if code == 0 {
		return false, errors.New("verify_signature returned 0")
	}
	return true, nil
}
