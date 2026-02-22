package task

import (
	"errors"
	"fmt"
)

const MaxLength = 100

var (
	ErrWrongLength = fmt.Sprintf("`title` must be less than or equal to %d characters", MaxLength)
	ErrEmptyTitle  = "`title` must not be empty"
)

type Title string

func NewTitle(title string) (*Title, error) {
	if len([]rune(title)) > MaxLength {
		return nil, errors.New(ErrWrongLength)
	}

	if len([]rune(title)) == 0 {
		return nil, errors.New(ErrEmptyTitle)
	}

	t := Title(title)
	return &t, nil
}

func (d Title) String() string {
	return string(d)
}
