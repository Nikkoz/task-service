package user

import (
	"fmt"
	"strings"
)

var (
	ErrEmptyPassword = fmt.Errorf("`PasswordHash` must not be empty")
)

type PasswordHash string

func NewPasswordHash(pass string) (*PasswordHash, error) {
	pass = strings.TrimSpace(pass)
	if len([]rune(pass)) == 0 {
		return nil, ErrEmptyPassword
	}

	p := PasswordHash(pass)
	return &p, nil
}

func (p PasswordHash) String() string {
	return string(p)
}
