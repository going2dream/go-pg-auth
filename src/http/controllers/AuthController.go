package AuthController

import (
	"encoding/json"
	l "github.com/ZeroDayDrake/go-pg-auth/src/logger"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	logger = l.New()
)

type (
	LoginReqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

func Login(ctx *fasthttp.RequestCtx) {
	var body LoginReqBody

	if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
		logger.Error(
			err.Error(),
			zap.ByteString("requestBody", ctx.PostBody()),
		)

		ctx.Error("Invalid data", fasthttp.StatusBadRequest)
		return
	}

	if !usernameValidate(body.Username, ctx) {
		return
	}

	//encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("hello"), bcrypt.DefaultCost)
	bcrypt.CompareHashAndPassword([]byte("JDJhJDEwJHpkTm14aHd0a3E1dUlodEpHSk8vMGU1bHdNLmNTMzdrNENiRXc0dlQwSWdLNDBXSjZrT09l"), []byte("hello"))

	ctx.SetStatusCode(fasthttp.StatusOK)

	response, _ := json.Marshal(map[string]interface{}{
		"token":   "encryptedPassword",
		"refresh": "hbhnjnknkkj",
	})

	if _, err := ctx.Write(response); err != nil {
		ctx.Error("Server Error", 500)
	}
}

func usernameValidate(username string, ctx *fasthttp.RequestCtx) bool {
	if len(username) < 3 {
		ctx.Error("Username length must be greater than 3 symbols", fasthttp.StatusUnprocessableEntity)
		return false
	}

	if len(username) > 255 {
		ctx.Error("Username length must be less than 255 symbols", fasthttp.StatusUnprocessableEntity)
		return false
	}

	return true
}

func RefreshToken(ctx *fasthttp.RequestCtx) {
	//ctx.SetStatusCode(fasthttp.StatusMovedPermanently)
	//ctx.Response.Header.Set("Location", "http://www.example.com/")
	ctx.WriteString("123")
}

func Logout(ctx *fasthttp.RequestCtx) {
	//ctx.SetStatusCode(fasthttp.StatusMovedPermanently)
	//ctx.Response.Header.Set("Location", "http://www.example.com/")
	ctx.WriteString("123")
}
