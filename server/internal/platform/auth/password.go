package auth

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(password string, hash string) bool
}

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{}
}

func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (h *BcryptPasswordHasher) Check(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// PasswordPolicy 密码策略接缝：客户项目可替换默认实现（如等保复杂度、历史密码检查）。
type PasswordPolicy interface {
	// Validate 校验候选密码，不符合策略时返回错误（建议为可直接展示给用户的文案）。
	Validate(password string) error
}

// DefaultPasswordPolicy 默认密码策略：至少 8 位，且同时包含字母和数字。
type DefaultPasswordPolicy struct{}

func (DefaultPasswordPolicy) Validate(password string) error {
	if len([]rune(password)) < 8 {
		return errors.New("密码长度至少 8 位")
	}
	var hasLetter, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("密码必须同时包含字母和数字")
	}
	return nil
}
