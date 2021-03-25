package utils

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/going2dream/go-pg-auth/src/app/logger"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"time"
)

var log = logger.New()

func JSONResponse(ctx *fasthttp.RequestCtx, response interface{}, code int) {
	ctx.SetStatusCode(code)
	ctx.SetContentType("application/json")
	if err := json.NewEncoder(ctx).Encode(response); err != nil {
		log.Error("JSON response encoder error", zap.String("details", err.Error()))
	}
}

var privateKey interface{}

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

type AppConfig struct {
	Environment          string        `yaml:"environment"`
	Domain               string        `yaml:"domain"`
	AuthURIPath          string        `yaml:"auth_uri_path"`
	BindIP               string        `yaml:"bind_ip"`
	BindPort             string        `yaml:"bind_port"`
	JWTTokenLifetime     time.Duration `yaml:"token_lifetime"`
	RefreshTokenLifetime time.Duration `yaml:"refresh_token_lifetime"`
}

var appConfig *AppConfig

func GetAppConfig() *AppConfig {
	if appConfig == nil {
		configFile, err := os.ReadFile("config/app.yml")
		if err != nil {
			log.Error("Can't open app config file", zap.String("details", err.Error()))
			return nil
		}

		type alias struct {
			Environment          string `yaml:"environment"`
			Domain               string `yaml:"domain"`
			AuthURIPath          string `yaml:"auth_uri_path"`
			BindIP               string `yaml:"bind_ip"`
			BindPort             string `yaml:"bind_port"`
			JWTTokenLifetime     string `yaml:"token_lifetime"`
			RefreshTokenLifetime string `yaml:"refresh_token_lifetime"`
		}

		var tmp alias
		if err := yaml.Unmarshal(configFile, &tmp); err != nil {
			log.Error("Unmarshal app config error", zap.Error(err))
		}

		t, err := time.ParseDuration(tmp.JWTTokenLifetime)
		if err != nil {
			log.Error("failed to parse to time.Duration", zap.Error(err))
		}

		r, err := time.ParseDuration(tmp.RefreshTokenLifetime)
		if err != nil {
			log.Error("failed to parse to time.Duration", zap.Error(err))
		}

		var config = AppConfig{
			Environment:          tmp.Environment,
			Domain:               tmp.Domain,
			AuthURIPath:          tmp.AuthURIPath,
			BindIP:               tmp.BindIP,
			BindPort:             tmp.BindPort,
			JWTTokenLifetime:     t,
			RefreshTokenLifetime: r,
		}

		appConfig = &config
	}

	return appConfig
}

func ClientIP(ctx *fasthttp.RequestCtx) string {
	clientIP := string(ctx.Request.Header.Peek("X-Forwarded-For"))
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
		//Get the first one, ie 1.1.1.1
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = strings.TrimSpace(string(ctx.Request.Header.Peek("X-Real-Ip")))
	if len(clientIP) > 0 {
		return clientIP
	}
	return ctx.RemoteIP().String()
}
