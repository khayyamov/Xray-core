#!/bin/bash
cat <<EOF > ./constant/constants.go
package constant

var (
	ENCRYPTED_CONFIG = false
	ENCRYPT_KEY      = "$ENCRYPT_KEY"
	ENCRYPT_KEY_IV   = "$ENCRYPT_KEY_IV"
)
EOF
