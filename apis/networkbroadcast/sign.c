#include "sign.h"
#include "secp256k1.h"
#include "secp256k1_recovery.h"

// https://github.com/bitcoin/bitcoin/blob/a6a860796a44a2805a58391a009ba22752f64e32/src/secp256k1/include/secp256k1_recovery.h#L81
void sign() {
	secp256k1_context* ctx = secp256k1_context_create(SECP256K1_CONTEXT_SIGN);
}
