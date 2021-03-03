package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func main() {
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
		Type:  "PRIVATE KEY",
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

	//privateKey := new(dsa.PrivateKey)
	//privateKey.PublicKey.Parameters = *params
	//
	//if err := dsa.GenerateKey(privateKey, rand.Reader); err != nil {
	//	panic(err)
	//}
	//
	//privateKeyFile, err := os.Create("keys/private.key")
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Save private key file
	//privateKeyEncoder := gob.NewEncoder(privateKeyFile)
	//
	//if err := privateKeyEncoder.Encode(privateKey); err != nil {
	//	panic(err)
	//}
	//
	//if err := privateKeyFile.Close(); err != nil {
	//	panic(err)
	//}

	//reader := rand.Reader
	//
	//var bitSize int
	//var outputDir string
	//flag.IntVar(&bitSize, "b", 2048, "bit size")
	//flag.StringVar(&outputDir, "output-dir", "../keys", "output dir")
	//flag.Parse()
	//
	//key, err := rsa.GenerateKey(reader, bitSize)
	//if err != nil {
	//	fmt.Println("Fatal error ", err.Error())
	//	os.Exit(1)
	//}
	//
	//publicKey := key.PublicKey
	//
	//keygen.SavePrivateKey(outputDir+"/private.key", key)
	//keygen.SavePublicKey(outputDir+"/public.pem", publicKey)
}
