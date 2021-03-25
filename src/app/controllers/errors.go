package controllers

type jsonT map[string]interface{}

var ErrBadCredentials = jsonT{
	"error": jsonT{
		"code":    1,
		"message": "These credentials do not match our records",
	},
}

var ErrInvalidRefreshToken = jsonT{
	"error": jsonT{
		"code":    2,
		"message": "Refresh token is invalid",
	},
}

var ErrInvalidRefreshSession = jsonT{
	"error": jsonT{
		"code":    3,
		"message": "Refresh session fingerprint is invalid",
	},
}

var ErrRefreshTokenExpired = jsonT{
	"error": jsonT{
		"code":    3,
		"message": "Refresh token has expired",
	},
}

var ErrRefreshTokenNotFound = jsonT{
	"error": jsonT{
		"code":    4,
		"message": "Refresh token not found",
	},
}
