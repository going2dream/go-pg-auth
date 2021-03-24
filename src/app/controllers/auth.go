package controllers

import (
	"encoding/json"
	"github.com/going2dream/go-pg-auth/src/app/logger"
	"github.com/going2dream/go-pg-auth/src/app/store"
	"github.com/going2dream/go-pg-auth/src/app/utils"
	"github.com/jackc/pgx/v4"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"time"
)

var log = logger.New()

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
		log.Error(
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
			utils.JSONResponse(ctx, ErrBadCredentials, 200)
			return
		}

		ctx.Error("Server error", fasthttp.StatusInternalServerError)
		return
	}

	if !user.ComparePassword(body.Password) {
		utils.JSONResponse(ctx, ErrBadCredentials, 200)
		return
	}

	// create a Square.jose DSA signer, used to sign the JWT
	var signerOpts = jose.SignerOptions{}
	signerOpts.WithType("JWT")

	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.EdDSA, Key: utils.GetPrivateKey()},
		&signerOpts,
	)
	if err != nil {
		log.Error("Failed to create signer", zap.String("details", err.Error()))
	}

	// create an instance of Builder that uses the dsa signer
	builder := jwt.Signed(signer)
	builder = builder.Claims(jwt.Claims{
		Subject: user.ID,
		Expiry:  jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
	})

	rawJWT, err := builder.CompactSerialize()
	if err != nil {
		log.Error("Failed to create JWT", zap.String("details", err.Error()))
		ctx.Error("Server error", 500)
		return
	}

	// Create refresh token

	utils.JSONResponse(
		ctx,
		map[string]interface{}{
			"token":   rawJWT,
			"refresh": "hbhnjnknkkj",
		},
		fasthttp.StatusOK,
	)
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
