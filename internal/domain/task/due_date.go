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

func (d *DueDate) DateTime() *time.Time {
	if d == nil {
		return nil
	}

	t := time.Time(*d)

	return &t
}

func (d *DueDate) String() string {
	if d == nil {
		return ""
	}

	return d.DateTime().String()
}
