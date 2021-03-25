package controllers

import (
	"encoding/json"
	"github.com/going2dream/go-pg-auth/src/app/logger"
	"github.com/going2dream/go-pg-auth/src/app/models"
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
		Login       string `json:"login"`
		Password    string `json:"password"`
		Fingerprint string `json:"fingerprint"`
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

	appConfig := utils.GetAppConfig()

	// create an instance of Builder that uses the dsa signer
	builder := jwt.Signed(signer)
	builder = builder.Claims(jwt.Claims{
		Subject: user.ID,
		Expiry:  jwt.NewNumericDate(time.Now().Add(appConfig.JWTTokenLifetime)),
	})

	rawJWT, err := builder.CompactSerialize()
	if err != nil {
		log.Error("Failed to create JWT", zap.String("details", err.Error()))
		ctx.Error("Server error", 500)
		return
	}

	// Create refresh token
	rt := &models.RefreshToken{
		UserID:      user.ID,
		UA:          string(ctx.Request.Header.UserAgent()),
		Fingerprint: body.Fingerprint,
		IP:          utils.ClientIP(ctx),
		ExpiresIn:   time.Now().Add(appConfig.RefreshTokenLifetime).Unix(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := rt.Validate(); err != nil {
		log.Error("Validation error", zap.Error(err))
		ctx.Error("Data is invalid", 400)
		return
	}

	if err := c.Store.RefreshToken().Create(rt); err != nil {
		log.Error("Refresh token creation error", zap.Error(err))
		ctx.Error("Server error", 500)
		return
	}

	rtCookie := fasthttp.Cookie{}
	rtCookie.SetKey("refreshToken")
	rtCookie.SetValue(rt.ID)
	rtCookie.SetExpire(time.Unix(rt.ExpiresIn, 0))
	rtCookie.SetDomain("." + appConfig.Domain)
	rtCookie.SetPath(appConfig.AuthURIPath)
	rtCookie.SetHTTPOnly(true)

	ctx.Response.Header.SetCookie(&rtCookie)

	utils.JSONResponse(
		ctx,
		map[string]interface{}{
			"token": rawJWT,
		},
		fasthttp.StatusOK,
	)
}

func (c *Auth) RefreshTokens(ctx *fasthttp.RequestCtx) {
	refreshTokenID := string(ctx.Request.Header.Cookie("refreshToken"))

	if refreshTokenID == "" {
		utils.JSONResponse(ctx, ErrInvalidRefreshToken, 200)
		return
	}

	requestBody := struct {
		Fingerprint string `json:"fingerprint"`
	}{}

	if err := json.Unmarshal(ctx.PostBody(), &requestBody); err != nil {
		log.Error(
			"Cant unmarshal refresh tokens request requestBody",
			zap.Error(err),
			zap.ByteString("requestBody", ctx.PostBody()),
		)

		ctx.Error("Invalid data", fasthttp.StatusBadRequest)
		return
	}

	refreshToken, err := c.Store.RefreshToken().Find(refreshTokenID)
	if err != nil {
		if err == pgx.ErrNoRows {
			utils.JSONResponse(ctx, ErrRefreshTokenNotFound, 200)
			ctx.Response.Header.DelCookie("refreshToken")
			return
		}

		log.Error(
			"Cant receive refresh token from database",
			zap.Error(err),
		)

		ctx.Error("Server error", fasthttp.StatusInternalServerError)
		return
	}

	if err := c.Store.RefreshToken().Delete(refreshToken.ID); err != nil {
		log.Error(
			"Cant delete refresh token from database",
			zap.Error(err),
		)

		ctx.Response.Header.DelCookie("refreshToken")
		ctx.Error("Server error", fasthttp.StatusInternalServerError)
		return
	}

	if !refreshToken.CompareFingerprint(requestBody.Fingerprint) {
		utils.JSONResponse(ctx, ErrInvalidRefreshSession, 200)
		ctx.Response.Header.DelCookie("refreshToken")

		log.Warn(
			"Unauthorized renewal of tokens",
			zap.String("tokenFingerprint", refreshToken.Fingerprint),
			zap.String("requestFingerprint", requestBody.Fingerprint),
			zap.String("requestIP", utils.ClientIP(ctx)),
		)

		return
	}

	if refreshToken.IsExpired() {
		utils.JSONResponse(ctx, ErrRefreshTokenExpired, 200)
		ctx.Response.Header.DelCookie("refreshToken")
		return
	}

	var signerOpts = jose.SignerOptions{}
	signerOpts.WithType("JWT")

	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.EdDSA, Key: utils.GetPrivateKey()},
		&signerOpts,
	)
	if err != nil {
		log.Error("Failed to create signer", zap.String("details", err.Error()))
	}

	appConfig := utils.GetAppConfig()

	// create an instance of Builder that uses the dsa signer
	builder := jwt.Signed(signer)
	builder = builder.Claims(jwt.Claims{
		Subject: refreshToken.UserID,
		Expiry:  jwt.NewNumericDate(time.Now().Add(appConfig.JWTTokenLifetime)),
	})

	rawJWT, err := builder.CompactSerialize()
	if err != nil {
		log.Error("Failed to create JWT", zap.String("details", err.Error()))
		ctx.Error("Server error", 500)
		return
	}

	// Recreate refresh token
	rt := &models.RefreshToken{
		UserID:      refreshToken.UserID,
		UA:          string(ctx.Request.Header.UserAgent()),
		Fingerprint: refreshToken.Fingerprint,
		IP:          utils.ClientIP(ctx),
		ExpiresIn:   time.Now().Add(appConfig.RefreshTokenLifetime).Unix(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := rt.Validate(); err != nil {
		log.Error("Validation error", zap.Error(err))
		ctx.Error("Server error", 500)
		return
	}

	if err := c.Store.RefreshToken().Create(rt); err != nil {
		log.Error("Refresh token creation error", zap.Error(err))
		ctx.Error("Server error", 500)
		return
	}

	rtCookie := fasthttp.Cookie{}
	rtCookie.SetKey("refreshToken")
	rtCookie.SetValue(rt.ID)
	rtCookie.SetExpire(time.Unix(rt.ExpiresIn, 0))
	rtCookie.SetDomain("." + appConfig.Domain)
	rtCookie.SetPath(appConfig.AuthURIPath)
	rtCookie.SetHTTPOnly(true)

	ctx.Response.Header.SetCookie(&rtCookie)

	utils.JSONResponse(
		ctx,
		map[string]interface{}{
			"token": rawJWT,
		},
		fasthttp.StatusOK,
	)
}

func (c *Auth) Logout(ctx *fasthttp.RequestCtx) {
	refreshTokenID := string(ctx.Request.Header.Cookie("refreshToken"))

	if refreshTokenID == "" {
		utils.JSONResponse(ctx, ErrInvalidRefreshToken, 200)
		return
	}

	requestBody := struct {
		Fingerprint string `json:"fingerprint"`
	}{}

	if err := json.Unmarshal(ctx.PostBody(), &requestBody); err != nil {
		log.Error(
			"Cant unmarshal refresh tokens request requestBody",
			zap.Error(err),
			zap.ByteString("requestBody", ctx.PostBody()),
		)

		ctx.Error("Invalid data", fasthttp.StatusBadRequest)
		return
	}

	refreshToken, err := c.Store.RefreshToken().Find(refreshTokenID)
	if err != nil {
		if err == pgx.ErrNoRows {
			utils.JSONResponse(ctx, ErrRefreshTokenNotFound, 200)
			ctx.Response.Header.DelCookie("refreshToken")
			return
		}

		log.Error(
			"Cant receive refresh token from database",
			zap.Error(err),
		)

		ctx.Error("Server error", fasthttp.StatusInternalServerError)
		return
	}

	if !refreshToken.CompareFingerprint(requestBody.Fingerprint) {
		utils.JSONResponse(ctx, ErrInvalidRefreshSession, 200)
		ctx.Response.Header.DelCookie("refreshToken")

		log.Warn(
			"Unauthorized logout",
			zap.String("tokenFingerprint", refreshToken.Fingerprint),
			zap.String("requestFingerprint", requestBody.Fingerprint),
			zap.String("requestIP", utils.ClientIP(ctx)),
		)

		return
	}

	if err := c.Store.RefreshToken().Delete(refreshTokenID); err != nil {
		log.Error(
			"Cant delete refresh token from database",
			zap.Error(err),
		)

		ctx.Error("Server error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.DelCookie("refreshToken")
	ctx.SetStatusCode(200)
}
