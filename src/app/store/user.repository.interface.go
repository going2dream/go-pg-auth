package store

import "github.com/ZeroDayDrake/go-pg-auth/src/app/models"

type UserRepository interface {
	//Create(*models.User) error
	Find(string) (*models.User, error)
	FindByLogin(string) (*models.User, error)
}
