package task

import (
	"errors"
	"strings"
)

type Description string

var ErrEmptyDescription = errors.New("`description` must not be empty")

func NewDescription(description string) (*Description, error) {
	description = strings.TrimSpace(description)
	if len([]rune(description)) == 0 {
		return nil, ErrEmptyDescription
	}

	d := Description(description)
	return &d, nil
}

func (d Description) String() string {
	return string(d)
}
