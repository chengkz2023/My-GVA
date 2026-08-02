package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
