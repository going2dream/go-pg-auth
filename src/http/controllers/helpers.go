package controllers

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"os"
)

var privateKey interface{}

func JSONResponse(ctx *fasthttp.RequestCtx, response interface{}, code int) {
	ctx.SetStatusCode(code)
	ctx.SetContentType("application/json")
	if err := json.NewEncoder(ctx).Encode(response); err != nil {
		log.Error("JSON response encoder error", zap.String("details", err.Error()))
	}
}

func GetPrivateKey() interface{} {
	if privateKey == nil {
		if _, err := os.Stat("keys/private.pem"); os.IsNotExist(err) {
			log.Error("Private key file is not exist")
			return nil
		}

		privateKeyFile, err := os.ReadFile("keys/private.pem")
		if err != nil {
			log.Error("Can't open private key file", zap.String("details", err.Error()))
			return nil
		}

		block, _ := pem.Decode([]byte(privateKeyFile))
		if block == nil {
			log.Error("Failed to parse PEM block containing the key")
			return nil
		}

		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			log.Error("Private key parse error", zap.String("details", err.Error()))
			return nil
		}

		privateKey = key
	}

	return privateKey
}
