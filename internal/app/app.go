package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nikkoz/task-service/internal/config"
	"github.com/Nikkoz/task-service/internal/repository/postgres"
	"github.com/Nikkoz/task-service/internal/service"
	"github.com/Nikkoz/task-service/internal/service/password"
	"github.com/Nikkoz/task-service/internal/service/token"
	"github.com/Nikkoz/task-service/internal/transport/http"
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/Nikkoz/task-service/pkg/logger"
)

func Run() {
	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("unable to parse environment variables: %v", err)
	}

	// create context
	ctx := context.Empty()
	defer ctx.Cancel()

	// create logger
	logger.New(cfg.App.Environment.IsProduction(), cfg.Log.Level.String())

	// create db
	db, err := NewDB(ctx, cfg.Db)
	if err != nil {
		_ = logger.ErrorWithContext(ctx, err)
	}
	defer db.Close()

	var (
		taskRepo    = postgres.NewTaskRepo(db)
		taskService = service.NewTaskService(taskRepo)

		userRepo       = postgres.NewUserRepo(db)
		passwordHasher = password.NewBcryptHasher(cfg.Auth.Cost)
		tokenManager   = token.NewJWTManager(cfg.Auth.JwtSecret, cfg.Auth.JwtTtl)
		authService    = service.NewAuthService(userRepo, passwordHasher, tokenManager)

		isProduction = cfg.App.Environment.IsProduction()
		listenerHttp = http.NewServer(taskService, authService, tokenManager, isProduction, http.Options{})
	)

	listenerHttp.Run(cfg.Http)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err := <-listenerHttp.Notify():
		logger.Error(fmt.Errorf("app - Run http server: %v", err))
	case done := <-ctx.Done():
		logger.Info(fmt.Sprintf("app - Run - ctx.Done: %v", done))
	}
}
