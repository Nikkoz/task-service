package auth

import (
	"github.com/Nikkoz/task-service/internal/domain/user"
	"github.com/Nikkoz/task-service/internal/service"
	"github.com/Nikkoz/task-service/pkg/context"
)

//go:generate mockery --name Service --output ./mocks --outpkg mocks
type Service interface {
	Register(ctx context.Context, in service.RegisterInput) (user.User, error)
	Login(ctx context.Context, in service.LoginInput) (service.LoginResult, error)
}
