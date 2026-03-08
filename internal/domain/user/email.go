package user

import (
	"fmt"
	"strings"
)

const MaxEmailLength = 100

var (
	ErrWrongEmailLength = fmt.Errorf("`email` must be less than or equal to %d characters", MaxEmailLength)
	ErrEmptyEmail       = fmt.Errorf("`email` must not be empty")
)

type Email string

func NewEmail(email string) (*Email, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	if len([]rune(email)) == 0 {
		return nil, ErrEmptyEmail
	}

	if len([]rune(email)) > MaxEmailLength {
		return nil, ErrWrongEmailLength
	}

	e := Email(email)
	return &e, nil
}

func (e Email) String() string {
	return string(e)
}
