#ifndef GOSTEEMRPC_SIGNING_H
#define GOSTEEMRPC_SIGNING_H

int sign_transaction(
	const unsigned char *digest,
	const unsigned char *privkey,
	char *signature,
	int *recid	
);

#endif
