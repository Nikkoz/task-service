package user

import (
	"fmt"
	"strings"
)

const MaxPassLength = 250

var (
	ErrWrongPassLength = fmt.Errorf("`PasswordHash` must be less than or equal to %d characters", MaxPassLength)
	ErrEmptyPassword   = fmt.Errorf("`PasswordHash` must not be empty")
)

type Password string

func NewPassword(pass string) (*Password, error) {
	pass = strings.TrimSpace(pass)
	if len([]rune(pass)) == 0 {
		return nil, ErrEmptyPassword
	}

	if len([]rune(pass)) > MaxPassLength {
		return nil, ErrWrongPassLength
	}

	p := Password(pass)
	return &p, nil
}

func (p Password) String() string {
	return string(p)
}
