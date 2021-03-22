package controllers

type jsonT map[string]interface{}

var ErrBadCredentials = jsonT{
	"error": jsonT{
		"code":    1,
		"message": "These credentials do not match our records",
	},
}
