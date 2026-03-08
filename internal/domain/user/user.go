package user

import "time"

type User struct {
	ID uint64

	// Email is a user's email.
	Email Email

	// PasswordHash is a user's hash of password.
	PasswordHash Password

	CreatedAt time.Time
	UpdatedAt time.Time
}
