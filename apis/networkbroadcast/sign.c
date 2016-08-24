#include "sign.h"
#include "secp256k1.h"
#include "secp256k1_recovery.h"
#include "ecdsa.h"

int sign()

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
