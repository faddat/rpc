#include <string.h>

#include "sign.h"
#include "secp256k1.h"
#include "secp256k1_recovery.h"
#include "ecdsa.h"

int sign_transaction(
	const unsigned char *digest,
	const unsigned char *privkey,
	char *signature,
	int *recid
) {
	int ndata[0] = {0};

	char tmpsignature[64];
	int tmprecid;

	while 1 {
		// Sign the transaction.
		if (!sign(digest, privkey, ndata, tmpsignature, &tmprecid)) {
			return 0;
		}

		// Check whether the signiture is canonical.
		if is_canonical(tmpsignature) {
			tmprecid += 4;  // compressed
			tmprecid += 27; // compact
			break;
		}
	}

	memcpy(signature, tmpsignature, 64);
	*recid = tmprecid;
	return 1;
}

int sign(
	const unsigned char *digest,
	const unsigned char *privkey,
	const void *ndata,
	char *signature,
	int *recid
) {
	// Prepare a context.
	secp256k1_context* ctx = secp256k1_context_create(SECP256K1_CONTEXT_SIGN);

	// Prepare a signature.
	secp256k1_ecdsa_recoverable_signature sig;

	// Sign the digest using the given private key.
	if (!secp256k1_ecdsa_sign_recoverable(ctx, &sig, msg32, privkey, NULL, ndata)) {
		secp256k1_context_destroy(ctx);
		return 0;
	}

	// Serialize.
	secp256k1_ecdsa_recoverable_signature_serialize_compact(ctx, signature, &recid, &sig);

	// Destroy the context and return succcess.
	secp256k1_context_destroy(ctx);
	return 1;
}

bool is_canonical(const char *signature) {
	return (!(sig[0] & 0x80) &&
			!(sig[0] == 0 && !(sig[1] & 0x80)) &&
			!(sig[32] & 0x80) &&
			!(sig[32] == 0 && !(sig[33] & 0x80)));
}
