package controllers

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/ZeroDayDrake/go-pg-auth/src/http/store"
	l "github.com/ZeroDayDrake/go-pg-auth/src/logger"
	"github.com/jackc/pgx/v4"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"os"
	"time"
)

var logger = l.New()

type (
	Auth struct {
		Store store.Store
	}

	LoginReqBody struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
)

func (c *Auth) Login(ctx *fasthttp.RequestCtx) {
	var body LoginReqBody

	if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
		logger.Error(
			err.Error(),
			zap.ByteString("requestBody", ctx.PostBody()),
		)

		ctx.Error("Invalid data", fasthttp.StatusBadRequest)
		return
	}

	if !c.usernameValidate(body.Login, ctx) {
		return
	}

	user, err := c.Store.User().FindByLogin(body.Login)
	if err != nil {
		if err == pgx.ErrNoRows {
			JSONError(ctx, ErrBadCredentials, 200)
			return
		}

		ctx.Error("Server error", fasthttp.StatusInternalServerError)
		return
	}

	if !user.ComparePassword(body.Password) {
		JSONError(ctx, ErrBadCredentials, 200)
		return
	}

	if _, err := os.Stat("keys/private.pem"); os.IsNotExist(err) {
		logger.Error("Private key file is not exist")
		ctx.Error("Server error", 500)
		return
	}

	privateKey, err := os.ReadFile("keys/private.pem")
	if err != nil {
		logger.Error("Can't open private key file", zap.String("details", err.Error()))
	}

	// create a Square.jose DSA signer, used to sign the JWT
	var signerOpts = jose.SignerOptions{}
	signerOpts.WithType("JWT")

	//var privateKeyInstance ed25519.PrivateKey = privateKey
	block, _ := pem.Decode(privateKey)
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Fatal("failed to decode PEM block containing public key")
	}

	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.EdDSA, Key: block},
		&signerOpts,
	)
	if err != nil {
		logger.Error("Failed to create signer", zap.String("details", err.Error()))
	}

	// create an instance of Builder that uses the dsa signer
	builder := jwt.Signed(signer)

	// public claims
	publicClaims := jwt.Claims{
		Issuer:   "issuer1",
		Subject:  "subject1",
		ID:       "id1",
		Audience: jwt.Audience{"aud1", "aud2"},
		IssuedAt: jwt.NewNumericDate(time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)),
		Expiry:   jwt.NewNumericDate(time.Date(2017, 1, 1, 0, 15, 0, 0, time.UTC)),
	}
	// private claims, as payload is JSON use the generic json patterns
	privateClaims := map[string]interface{}{
		"privateClaim1": "val1",
		"privateClaim2": []string{"val2", "val3"},
		"anyJSONObjectClaim": map[string]interface{}{
			"name": "john",
			"phones": map[string]string{
				"phone1": "123",
				"phone2": "456",
			},
		},
	}
	builder = builder.Claims(publicClaims).Claims(privateClaims)

	// validate all ok, sign with the RSA key, and return a compact JWT
	rawJWT, err := builder.CompactSerialize()
	if err != nil {
		logger.Error("Failed to create JWT", zap.String("details", err.Error()))
		ctx.Error("Server error", 500)
		return
	}

	fmt.Println(rawJWT)

	//encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	//fmt.Println(string(encryptedPassword))

	//bcrypt.CompareHashAndPassword([]byte("JDJhJDEwJHpkTm14aHd0a3E1dUlodEpHSk8vMGU1bHdNLmNTMzdrNENiRXc0dlQwSWdLNDBXSjZrT09l"), []byte("hello"))

	ctx.SetStatusCode(fasthttp.StatusOK)

	response, _ := json.Marshal(map[string]interface{}{
		"token":   "encryptedPassword",
		"refresh": "hbhnjnknkkj",
	})

	if _, err := ctx.Write(response); err != nil {
		ctx.Error("Server Error", 500)
	}
}

func (c *Auth) usernameValidate(username string, ctx *fasthttp.RequestCtx) bool {
	if len(username) < 3 {
		ctx.Error("DBUsername length must be greater than 3 symbols", fasthttp.StatusUnprocessableEntity)
		return false
	}

	if len(username) > 255 {
		ctx.Error("DBUsername length must be less than 255 symbols", fasthttp.StatusUnprocessableEntity)
		return false
	}

	return true
}

func (c *Auth) RefreshToken(ctx *fasthttp.RequestCtx) {
	//ctx.SetStatusCode(fasthttp.StatusMovedPermanently)
	//ctx.Response.Header.Set("Location", "http://www.example.com/")
	ctx.WriteString("123")
}

func (c *Auth) Logout(ctx *fasthttp.RequestCtx) {
	//ctx.SetStatusCode(fasthttp.StatusMovedPermanently)
	//ctx.Response.Header.Set("Location", "http://www.example.com/")
	ctx.WriteString("123")
}
