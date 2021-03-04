package models

import "golang.org/x/crypto/bcrypt"

type RefreshToken struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}

// Validate ...
//func (u *User) Validate() error {
//	return validation.ValidateStruct(
//		u,
//		validation.Field(&u.Email, validation.Required, is.Email),
//		validation.Field(&u.DBPassword, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
//	)
//}

// BeforeCreate ...
func (u *RefreshToken) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.Password = enc
	}

	return nil
}

// Sanitize ...
func (u *RefreshToken) Sanitize() {
	u.Password = ""
}

// ComparePassword ...
func (u *RefreshToken) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
