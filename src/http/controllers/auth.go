package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/ZeroDayDrake/go-pg-auth/src/http/store"
	l "github.com/ZeroDayDrake/go-pg-auth/src/logger"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	logger = l.New()
)

type (
	Auth struct {
		store *store.Store
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

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	fmt.Println(string(encryptedPassword))

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
