#include <stdbool.h>
#include <stdio.h>
#include <string.h>

#include "secp256k1.h"
#include "secp256k1_recovery.h"

#include "signing.h"

static int sign(
	const secp256k1_context* ctx,
	const unsigned char *digest,
	const unsigned char *privkey,
	const void *ndata,
	unsigned char *signature,
	int *recid
);

static bool is_canonical(const unsigned char *signature);

void dump(const unsigned char *array, int len) {
	for (int i = 0; i < len; i++) {
		printf("%d ", array[i]);
	}
	printf("\n");
}

int sign_transaction(
	const unsigned char *digest,
	const unsigned char *privkey,
	unsigned char *signature,
	int *recid
) {
	secp256k1_context* ctx = secp256k1_context_create(SECP256K1_CONTEXT_SIGN);

	int ndata = 1;

	while (1) {
		// Sign the transaction.
		if (!sign(ctx, digest, privkey, &ndata, signature, recid)) {
			secp256k1_context_destroy(ctx);
			return 0;
		}

		// Check whether the signiture is canonical.
		if (is_canonical(signature)) {
			*recid += 4;  // compressed
			*recid += 27; // compact
			break;
		}

		ndata++;
	}

	secp256k1_context_destroy(ctx);
	return 1;
}

static int sign(
	const secp256k1_context* ctx,
	const unsigned char *digest,
	const unsigned char *privkey,
	const void *ndata,
	unsigned char *signature,
	int *recid
) {
	printf("DIGEST:\n");
	dump(digest, 32);
	printf("KEY:\n");
	dump(privkey, 32);

	// Prepare a signature.
	secp256k1_ecdsa_recoverable_signature sig;

	// Sign the digest using the given private key.
	if (!secp256k1_ecdsa_sign_recoverable(ctx, &sig, digest, privkey, NULL, ndata)) {
		return 0;
	}

	printf("SIGNATURE DATA:\n");
	dump(sig.data, 65);

	// Serialize and return success.
	secp256k1_ecdsa_recoverable_signature_serialize_compact(ctx, signature, recid, &sig);
	return 1;
}

static bool is_canonical(const unsigned char *sig) {
	return (!(sig[0] & 0x80) &&
	        !(sig[0] == 0 && !(sig[1] & 0x80)) &&
	        !(sig[32] & 0x80) &&
	        !(sig[32] == 0 && !(sig[33] & 0x80)));
}
