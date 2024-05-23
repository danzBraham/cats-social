package usecases

import (
	"context"
	"time"

	"github.com/danzbraham/cats-social/internal/applications/securities"
	user_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/users"

	user_entity "github.com/danzbraham/cats-social/internal/domains/entities/users"
	"github.com/danzbraham/cats-social/internal/domains/repositories"
	"github.com/oklog/ulid/v2"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, payload *user_entity.RegisterUserRequest) (*user_entity.RegisterUserResponse, error)
}

type UserUsecaseImpl struct {
	Repository       repositories.UserRepository
	PasswordHasher   securities.PasswordHasher
	AuthTokenManager securities.AuthTokenManager
}

func NewUserUsecase(repository repositories.UserRepository, passwordHasher securities.PasswordHasher, authTokenManager securities.AuthTokenManager) UserUsecase {
	return &UserUsecaseImpl{
		Repository:       repository,
		PasswordHasher:   passwordHasher,
		AuthTokenManager: authTokenManager,
	}
}

func (uc *UserUsecaseImpl) RegisterUser(ctx context.Context, payload *user_entity.RegisterUserRequest) (*user_entity.RegisterUserResponse, error) {
	isEmailExists, err := uc.Repository.VerifyEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}
	if isEmailExists {
		return nil, user_exception.ErrEmailAlreadyExists
	}

	id := ulid.Make().String()

	hashedPassword, err := uc.PasswordHasher.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	user := &user_entity.User{
		ID:       id,
		Email:    payload.Email,
		Name:     payload.Name,
		Password: hashedPassword,
	}

	if err := uc.Repository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	accessToken, err := uc.AuthTokenManager.GenerateToken(2*time.Hour, id)
	if err != nil {
		return nil, err
	}

	res := &user_entity.RegisterUserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Email,
		AccessToken: accessToken,
	}

	return res, nil
}
