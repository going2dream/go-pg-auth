package store

import "github.com/ZeroDayDrake/go-pg-auth/src/http/models"

type RefreshTokenRepository interface {
	Create(*models.RefreshToken) error
	Find(string) (*models.RefreshToken, error)
	Delete(*models.RefreshToken) error
}
