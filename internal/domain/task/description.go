package task

import "errors"

type Description string

var ErrEmptyDescription = "`title` must not be empty"

func NewDescription(description string) (*Description, error) {
	if len([]rune(description)) == 0 {
		return nil, errors.New(ErrEmptyDescription)
	}

	d := Description(description)
	return &d, nil
}

func (d Description) String() string {
	return string(d)
}
