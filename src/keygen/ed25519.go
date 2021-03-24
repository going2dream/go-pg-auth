package keygen

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func MakeEd25519Keys() {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)

	// Save private key
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	var privateKeyPEM = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privateKeyFile, err := os.Create("keys/private.pem")
	if err != nil {
		panic(err)
	}

	err = pem.Encode(privateKeyFile, privateKeyPEM)
	if err != nil {
		panic(err)
	}

	if err := privateKeyFile.Close(); err != nil {
		panic(err)
	}

	// Save public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	var publicKeyPEM = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicKeyFile, err := os.Create("keys/public.pem")
	if err != nil {
		panic(err)
	}

	err = pem.Encode(publicKeyFile, publicKeyPEM)
	if err != nil {
		panic(err)
	}

	if err := publicKeyFile.Close(); err != nil {
		panic(err)
	}
}
