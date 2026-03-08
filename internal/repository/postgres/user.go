package postgres

import (
	"errors"

	"github.com/Nikkoz/task-service/internal/domain/user"
	"github.com/Nikkoz/task-service/internal/repository"
	"github.com/Nikkoz/task-service/internal/service"
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/Nikkoz/task-service/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ service.UserRepository = (*UserRepo)(nil)

type UserRepo struct {
	db DBTX
}

func NewUserRepo(db DBTX) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u user.User) (user.User, error) {
	const q = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, created_at, updated_at
	`

	var out user.User
	var rawEmail string

	err := r.db.QueryRow(ctx, q,
		u.Email.String(),
		u.PasswordHash,
	).Scan(
		&out.ID,
		&rawEmail,
		&out.PasswordHash,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return user.User{}, repository.ErrAlreadyExists
		}

		return user.User{}, logger.ErrorWithContext(ctx, err)
	}

	email, _ := user.NewEmail(rawEmail)
	out.Email = *email

	return out, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uint64) (user.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var out user.User
	var rawEmail string

	err := r.db.QueryRow(ctx, q, id).Scan(
		&out.ID,
		&rawEmail,
		&out.PasswordHash,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, repository.ErrNotFound
		}

		return user.User{}, logger.ErrorWithContext(ctx, err)
	}

	email, _ := user.NewEmail(rawEmail)
	out.Email = *email

	return out, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (user.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
	`

	var out user.User
	var rawEmail string

	err := r.db.QueryRow(ctx, q, email).Scan(
		&out.ID,
		&rawEmail,
		&out.PasswordHash,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, repository.ErrNotFound
		}

		return user.User{}, logger.ErrorWithContext(ctx, err)
	}

	normalizedEmail, _ := user.NewEmail(rawEmail)
	out.Email = *normalizedEmail

	return out, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
