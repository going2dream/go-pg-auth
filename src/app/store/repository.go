package store

import "github.com/going2dream/go-pg-auth/src/app/models"

type UserRepository interface {
	//Create(*models.User) error
	Find(string) (*models.User, error)
	FindByLogin(string) (*models.User, error)
}

type RefreshTokenRepository interface {
	Create(*models.RefreshToken) error
	Find(string) (*models.RefreshToken, error)
	Delete(*models.RefreshToken) error
}
