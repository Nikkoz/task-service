package task

import (
	"fmt"
	"strings"
)

const MaxLength = 100

var (
	ErrWrongLength = fmt.Errorf("`title` must be less than or equal to %d characters", MaxLength)
	ErrEmptyTitle  = fmt.Errorf("`title` must not be empty")
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

func (t Title) String() string {
	return string(t)
}
