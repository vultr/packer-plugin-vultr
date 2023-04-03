// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func PublicKeyFromPrivate(privateKeyBytes []byte) ([]byte, error) {
	key, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("Error on parsing SSH private key: %s", err)
	}

	return ssh.MarshalAuthorizedKey(key.PublicKey()), nil
}
