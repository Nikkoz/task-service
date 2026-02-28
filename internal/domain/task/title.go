package task

import (
	"errors"
	"fmt"
	"strings"
)

const MaxLength = 100

var (
	ErrWrongLength = errors.New(fmt.Sprintf("`title` must be less than or equal to %d characters", MaxLength))
	ErrEmptyTitle  = errors.New("`title` must not be empty")
)

type Title string

func NewTitle(title string) (*Title, error) {
	title = strings.TrimSpace(title)
	if len([]rune(title)) == 0 {
		return nil, ErrEmptyTitle
	}

	if len([]rune(title)) > MaxLength {
		return nil, ErrWrongLength
	}

	t := Title(title)
	return &t, nil
}

func (d Title) String() string {
	return string(d)
}
