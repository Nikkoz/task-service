package postgres

import (
	"errors"
	"os"
	"testing"

	"github.com/Nikkoz/task-service/internal/testutil"
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"

	"github.com/Nikkoz/task-service/internal/domain/user"
	"github.com/Nikkoz/task-service/internal/repository"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("ENV_FILE", ".testing")

	code := m.Run()

	testutil.ClosePool()
	os.Exit(code)
}

func TestUserRepo_CreateAndGetByID(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		repo := NewUserRepo(tx)
		assertion := assert.New(t)

		email, _ := user.NewEmail("test@example.com")
		password, _ := user.NewPassword("password")

		created, err := repo.Create(ctx, user.User{
			Email:        *email,
			PasswordHash: *password,
		})
		assertion.NoError(err)

		found, err := repo.GetByID(ctx, created.ID)
		assertion.NoError(err)

		assertion.Equal(created.ID, found.ID)
		assertion.Equal(created.Email.String(), found.Email.String())
	})
}

func TestUserRepo_GetByEmail(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		repo := NewUserRepo(tx)
		assertion := assert.New(t)

		email, _ := user.NewEmail("test@example.com")
		password, _ := user.NewPassword("password")

		created, err := repo.Create(ctx, user.User{
			Email:        *email,
			PasswordHash: *password,
		})
		assertion.NoError(err)

		found, err := repo.GetByEmail(ctx, email.String())
		assertion.NoError(err)

		assertion.Equal(created.ID, found.ID)
	})
}

func TestUserRepo_Create_DuplicateEmail(t *testing.T) {
	testutil.WithTx(t, func(ctx context.Context, tx pgx.Tx) {
		repo := NewUserRepo(tx)
		assertion := assert.New(t)

		email, _ := user.NewEmail("test@example.com")
		password, _ := user.NewPassword("password")

		_, err := repo.Create(ctx, user.User{
			Email:        *email,
			PasswordHash: *password,
		})
		assertion.NoError(err)

		_, err = repo.Create(ctx, user.User{
			Email:        *email,
			PasswordHash: "hash",
		})
		assertion.Error(err)

		assertion.True(errors.Is(err, repository.ErrAlreadyExists))
	})
}
