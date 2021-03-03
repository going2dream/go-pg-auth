package controllers

import (
	"encoding/json"
	l "github.com/ZeroDayDrake/go-pg-auth/src/logger"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var log = l.New()

type jsonType map[string]interface{}

func JSONError(ctx *fasthttp.RequestCtx, err interface{}, code int) {
	ctx.SetStatusCode(code)
	ctx.SetContentType("application/json")
	if err := json.NewEncoder(ctx).Encode(err); err != nil {
		log.Error("JSON error encoder error", zap.String("details", err.Error()))
	}
}

var ErrBadCredentials = jsonType{
	"error": jsonType{
		"code":    1,
		"message": "These credentials do not match our records",
	},
}
