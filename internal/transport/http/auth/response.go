package auth

import (
	"github.com/Nikkoz/task-service/internal/domain/user"
	"github.com/Nikkoz/task-service/internal/service"
)

type (
	registerResponse struct {
		ID    uint64 `json:"id"`
		Email string `json:"email"`
	}

	loginResponse struct {
		Token string `json:"token"`
	}
)

func toRegisterResponse(user user.User) *registerResponse {
	return &registerResponse{
		ID:    user.ID,
		Email: user.Email.String(),
	}
}

func toLoginResponse(result service.LoginResult) *loginResponse {
	return &loginResponse{
		Token: result.Token,
	}
}
