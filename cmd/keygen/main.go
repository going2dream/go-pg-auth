package main

import (
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	"github.com/ZeroDayDrake/go-pg-auth/src/keygen"
	"os"
)

func main() {
	reader := rand.Reader

	var bitSize int
	var outputDir string
	flag.IntVar(&bitSize, "b", 2048, "bit size")
	flag.StringVar(&outputDir, "output-dir", "../keys", "output dir")
	flag.Parse()

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}

	publicKey := key.PublicKey

	keygen.SavePEMKey(outputDir+"/private.pem", key)
	keygen.SavePublicPEMKey(outputDir+"/public.pem", publicKey)
}
