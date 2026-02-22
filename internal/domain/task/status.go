package task

import "fmt"

type Status string

const (
	StatusPlanned    Status = "planned"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

func NewStatus(status string) (*Status, error) {
	s := Status(status)
	if !s.Valid() {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	return &s, nil
}

func (s Status) Valid() bool {
	switch s {
	case StatusPlanned, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (s Status) String() string {
	return string(s)
}
