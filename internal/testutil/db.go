package testutil

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	poolOnce sync.Once
	poolVal  *pgxpool.Pool
	poolErr  error
)

func GetPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	poolOnce.Do(func() {
		config := GetConfig(t).Db
		sslMode := "enable"
		if !config.SslMode {
			sslMode = "disable"
		}

		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.Name,
			sslMode,
		)

		poolVal, poolErr = pgxpool.New(context.Empty(), dsn)
	})

	if poolErr != nil {
		t.Fatalf("app.NewDB: %v", poolErr)
	}

	return poolVal
}

// ClosePool closes pool if it was created.
// Call it from TestMain after m.Run().
func ClosePool() {
	if poolVal != nil {
		poolVal.Close()
	}
}

// WithTx runs fn inside a DB transaction and always rollbacks after the test.
// No truncates, no cleanup in tests.
func WithTx(t *testing.T, fn func(ctx context.Context, tx pgx.Tx)) {
	t.Helper()

	// небольшой timeout, чтобы тесты не зависали
	ctx := context.NewWithTimeout(context.Empty(), 205*time.Second)
	t.Cleanup(ctx.Cancel)

	tx, err := GetPool(t).Begin(ctx)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}

	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	fn(ctx, tx)
}
