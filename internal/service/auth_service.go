package service

import (
	"errors"

	"github.com/Nikkoz/task-service/internal/domain/user"
	"github.com/Nikkoz/task-service/internal/repository"
	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/Nikkoz/task-service/pkg/logger"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	users  UserRepository
	hasher PasswordHasher
	tokens TokenManager
}

func NewAuthService(users UserRepository, hasher PasswordHasher, tokens TokenManager) *AuthService {
	return &AuthService{
		users:  users,
		hasher: hasher,
		tokens: tokens,
	}
}

type RegisterInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token string
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (user.User, error) {
	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return user.User{}, logger.ErrorWithContext(ctx, err)
	}

	email, err := user.NewEmail(in.Email)
	if err != nil {
		return user.User{}, logger.ErrorWithContext(ctx, err)
	}

	passHash, err := user.NewPasswordHash(hash)
	if err != nil {
		return user.User{}, logger.ErrorWithContext(ctx, err)
	}

	created, err := s.users.Create(ctx, user.User{
		Email:        *email,
		PasswordHash: *passHash,
	})
	if err != nil {
		return user.User{}, logger.ErrorWithContext(ctx, err)
	}

	return created, nil
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (LoginResult, error) {
	found, err := s.users.GetByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return LoginResult{}, ErrInvalidCredentials
		}

		return LoginResult{}, logger.ErrorWithContext(ctx, err)
	}

	if err := s.hasher.Compare(found.PasswordHash.String(), in.Password); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(found.ID)
	if err != nil {
		return LoginResult{}, logger.ErrorWithContext(ctx, err)
	}

	return LoginResult{Token: token}, nil
}
