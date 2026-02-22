package task

import "time"

type Task struct {
	ID uint64

	// Title is a short title (max 100 chars in DB).
	Title Title

	// Description is an optional description (stored as text).
	Description Description

	// Status is the lifecycle state: planned | in_progress | done.
	Status Status

	// DueDate is an optional deadline in UTC (nil means no deadline).
	DueDate *DueDate

	CreatedAt time.Time
	UpdatedAt time.Time
}
