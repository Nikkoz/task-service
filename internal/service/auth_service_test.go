package service

import (
	"errors"
	"testing"

	"github.com/Nikkoz/task-service/pkg/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Nikkoz/task-service/internal/domain/user"
	"github.com/Nikkoz/task-service/internal/repository"
	"github.com/Nikkoz/task-service/internal/service/mocks"
)

func TestAuthService_Register_Success(t *testing.T) {
	users := mocks.NewUserRepository(t)
	hasher := mocks.NewPasswordHasher(t)
	tokens := mocks.NewTokenManager(t)

	svc := NewAuthService(users, hasher, tokens)

	email := "test@example.com"
	password := "secret123"
	in := RegisterInput{
		Email:    email,
		Password: password,
	}

	hashed := "hashed-secret"
	hasher.
		On("Hash", password).
		Return(hashed, nil).
		Once()

	login, _ := user.NewEmail(email)
	pass, _ := user.NewPasswordHash(hashed)
	users.
		On("Create", mock.Anything, mock.MatchedBy(func(u user.User) bool {
			return u.Email.String() == email &&
				u.PasswordHash.String() == hashed
		})).
		Return(user.User{
			ID:           1,
			Email:        *login,
			PasswordHash: *pass,
		}, nil).
		Once()

	out, err := svc.Register(context.Empty(), in)

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.Equal(uint64(1), out.ID)
	assertion.Equal(email, out.Email.String())

	hasher.AssertExpectations(t)
	users.AssertExpectations(t)
	tokens.AssertExpectations(t)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	users := mocks.NewUserRepository(t)
	hasher := mocks.NewPasswordHasher(t)
	tokens := mocks.NewTokenManager(t)

	svc := NewAuthService(users, hasher, tokens)

	email := "test@example.com"
	password := "secret123"
	in := RegisterInput{
		Email:    email,
		Password: password,
	}

	hashed := "hashed-secret"
	hasher.
		On("Hash", password).
		Return(hashed, nil).
		Once()

	users.
		On("Create", mock.Anything, mock.Anything).
		Return(user.User{}, repository.ErrAlreadyExists).
		Once()

	_, err := svc.Register(context.Empty(), in)

	assertion := assert.New(t)
	assertion.Error(err)
	assertion.ErrorAs(err, &repository.ErrAlreadyExists)

	hasher.AssertExpectations(t)
	users.AssertExpectations(t)
	tokens.AssertExpectations(t)
}

func TestAuthService_Register_HashError(t *testing.T) {
	users := mocks.NewUserRepository(t)
	hasher := mocks.NewPasswordHasher(t)
	tokens := mocks.NewTokenManager(t)

	svc := NewAuthService(users, hasher, tokens)

	email := "test@example.com"
	password := "secret123"
	in := RegisterInput{
		Email:    email,
		Password: password,
	}

	hasher.
		On("Hash", password).
		Return("", errors.New("hash failed")).
		Once()

	_, err := svc.Register(context.Empty(), in)

	assertion := assert.New(t)
	assertion.Error(err)

	hasher.AssertExpectations(t)
	users.AssertExpectations(t)
	tokens.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	users := mocks.NewUserRepository(t)
	hasher := mocks.NewPasswordHasher(t)
	tokens := mocks.NewTokenManager(t)

	svc := NewAuthService(users, hasher, tokens)

	email := "test@example.com"
	password := "secret123"
	hashed := "hashed-secret"
	in := LoginInput{
		Email:    email,
		Password: password,
	}

	id := uint64(42)
	login, _ := user.NewEmail(email)
	pass, _ := user.NewPasswordHash(hashed)

	users.
		On("GetByEmail", mock.Anything, email).
		Return(user.User{
			ID:           id,
			Email:        *login,
			PasswordHash: *pass,
		}, nil).
		Once()

	hasher.
		On("Compare", hashed, password).
		Return(nil).
		Once()

	token := "jwt-token"
	tokens.
		On("Generate", id).
		Return(token, nil).
		Once()

	out, err := svc.Login(context.Empty(), in)

	assertion := assert.New(t)
	assertion.NoError(err)
	assertion.Equal(token, out.Token)

	users.AssertExpectations(t)
	hasher.AssertExpectations(t)
	tokens.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	users := mocks.NewUserRepository(t)
	hasher := mocks.NewPasswordHasher(t)
	tokens := mocks.NewTokenManager(t)

	svc := NewAuthService(users, hasher, tokens)

	email := "test@example.com"
	password := "secret123"
	in := LoginInput{
		Email:    email,
		Password: password,
	}

	users.
		On("GetByEmail", mock.Anything, email).
		Return(user.User{}, repository.ErrNotFound).
		Once()

	_, err := svc.Login(context.Empty(), in)

	assertion := assert.New(t)
	assertion.Error(err)
	assertion.ErrorAs(err, &repository.ErrNotFound)

	users.AssertExpectations(t)
	hasher.AssertExpectations(t)
	tokens.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	users := mocks.NewUserRepository(t)
	hasher := mocks.NewPasswordHasher(t)
	tokens := mocks.NewTokenManager(t)

	svc := NewAuthService(users, hasher, tokens)

	email := "test@example.com"
	password := "wrong-password"
	in := LoginInput{
		Email:    email,
		Password: password,
	}

	hashed := "hashed-secret"
	login, _ := user.NewEmail(email)
	pass, _ := user.NewPasswordHash(hashed)
	users.
		On("GetByEmail", mock.Anything, email).
		Return(user.User{
			ID:           42,
			Email:        *login,
			PasswordHash: *pass,
		}, nil).
		Once()

	hasher.
		On("Compare", hashed, password).
		Return(errors.New("password mismatch")).
		Once()

	_, err := svc.Login(context.Empty(), in)

	assertion := assert.New(t)
	assertion.Error(err)
	assertion.ErrorAs(err, &ErrInvalidCredentials)

	users.AssertExpectations(t)
	hasher.AssertExpectations(t)
	tokens.AssertExpectations(t)
}

func TestAuthService_Login_TokenError(t *testing.T) {
	users := mocks.NewUserRepository(t)
	hasher := mocks.NewPasswordHasher(t)
	tokens := mocks.NewTokenManager(t)

	svc := NewAuthService(users, hasher, tokens)

	email := "test@example.com"
	password := "secret123"
	in := LoginInput{
		Email:    email,
		Password: password,
	}

	hashed := "hashed-secret"
	id := uint64(42)
	login, _ := user.NewEmail(email)
	pass, _ := user.NewPasswordHash(hashed)
	users.
		On("GetByEmail", mock.Anything, email).
		Return(user.User{
			ID:           id,
			Email:        *login,
			PasswordHash: *pass,
		}, nil).
		Once()

	hasher.
		On("Compare", hashed, password).
		Return(nil).
		Once()

	tokens.
		On("Generate", id).
		Return("", errors.New("token generation failed")).
		Once()

	_, err := svc.Login(context.Empty(), in)

	assertion := assert.New(t)
	assertion.Error(err)

	users.AssertExpectations(t)
	hasher.AssertExpectations(t)
	tokens.AssertExpectations(t)
}
