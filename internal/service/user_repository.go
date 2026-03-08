package service

import (
	"github.com/Nikkoz/task-service/internal/domain/user"
	"github.com/Nikkoz/task-service/pkg/context"
)

//go:generate mockery --name UserRepository --output ./mocks --outpkg mocks
type UserRepository interface {
	Create(ctx context.Context, u user.User) (user.User, error)
	GetByID(ctx context.Context, id uint64) (user.User, error)
	GetByEmail(ctx context.Context, email string) (user.User, error)
}
