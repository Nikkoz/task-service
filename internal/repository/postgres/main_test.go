//go:build integration

package postgres

import (
	"os"
	"testing"

	"github.com/Nikkoz/task-service/internal/domain/user"
	"github.com/Nikkoz/task-service/internal/testutil"
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/jackc/pgx/v5"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("ENV_FILE", ".testing")

	code := m.Run()

	testutil.ClosePool()
	os.Exit(code)
}

func createUser(t *testing.T, ctx context.Context, tx pgx.Tx, email string) (user.User, error) {
	t.Helper()

	repo := NewUserRepo(tx)

	e, _ := user.NewEmail(email)
	password, _ := user.NewPasswordHash("password")

	return repo.Create(ctx, user.User{
		Email:        *e,
		PasswordHash: *password,
	})
}
