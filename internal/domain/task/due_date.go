package task

import (
	"errors"
	"time"
)

type DueDate time.Time

var ErrDateInThePast = errors.New("`DueDate` must be in the future")

func NewDueDate(dueDate time.Time) (*DueDate, error) {
	if dueDate.Before(time.Now()) || dueDate.Equal(time.Now()) {
		return nil, ErrDateInThePast
	}

	d := DueDate(dueDate)
	return &d, nil
}

func (d DueDate) DateTime() time.Time {
	return time.Time(d)
}

func (d DueDate) String() string {
	return d.DateTime().String()
}
