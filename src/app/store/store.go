package store

type Store interface {
	User() UserRepository
	RefreshToken() RefreshTokenRepository
}
