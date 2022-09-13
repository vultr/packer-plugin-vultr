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
