package models

import (
	"errors"
	"time"
)

type RefreshToken struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	UA          string    `json:"ua"`
	Fingerprint string    `json:"fingerprint"`
	IP          string    `json:"ip"`
	ExpiresIn   int64     `json:"expires_in"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *RefreshToken) Validate() error {
	if t.UserID == "" {
		return errors.New("user id is required")
	}

	if t.UA == "" {
		return errors.New("user agent is required")
	}

	if t.Fingerprint == "" {
		return errors.New("fingerprint is required")
	}

	if t.IP == "" {
		return errors.New("ip is required")
	}

	if t.ExpiresIn == 0 {
		return errors.New("expires_in is required")
	}

	return nil
}

func (t *RefreshToken) CompareFingerprint(fingerprint string) bool {
	return fingerprint == t.Fingerprint
}

func (t *RefreshToken) IsExpired() bool {
	return time.Now().Unix() > t.ExpiresIn
}

//func (u *RefreshTokens) BeforeCreate() error {
//	if len(u.Password) > 0 {
//		enc, err := encryptString(u.Password)
//		if err != nil {
//			return err
//		}
//
//		u.Password = enc
//	}
//
//	return nil
//}

// Sanitize ...
//func (u *RefreshTokens) Sanitize() {
//	u.Password = ""
//}

// ComparePassword ...
//func (u *RefreshTokens) ComparePassword(password string) bool {
//	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
//}
