package services

import (
	"context"
	"time"

	auth_token_manager "github.com/danzbraham/cats-social/internal/commons/auth-token-manager"
	user_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/users"
	password_hasher "github.com/danzbraham/cats-social/internal/commons/password-hasher"
	user_entity "github.com/danzbraham/cats-social/internal/entities/user"

	"github.com/danzbraham/cats-social/internal/repositories"

	"github.com/oklog/ulid/v2"
)

type UserService interface {
	RegisterUser(ctx context.Context, payload *user_entity.RegisterUserRequest) (*user_entity.RegisterUserResponse, error)
	LoginUser(ctx context.Context, payload *user_entity.LoginUserRequest) (*user_entity.LoginUserResponse, error)
}

type UserServiceImpl struct {
	Repository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) UserService {
	return &UserServiceImpl{Repository: repository}
}

func (s *UserServiceImpl) RegisterUser(ctx context.Context, payload *user_entity.RegisterUserRequest) (*user_entity.RegisterUserResponse, error) {
	isEmailExists, err := s.Repository.VerifyEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}
	if isEmailExists {
		return nil, user_exception.ErrEmailAlreadyExists
	}

	hashedPassword, err := password_hasher.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	user := &user_entity.User{
		Id:       ulid.Make().String(),
		Email:    payload.Email,
		Name:     payload.Name,
		Password: hashedPassword,
	}

	if err := s.Repository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	accessToken, err := auth_token_manager.GenerateToken(2*time.Hour, user.Id)
	if err != nil {
		return nil, err
	}

	return &user_entity.RegisterUserResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil
}

func (s *UserServiceImpl) LoginUser(ctx context.Context, payload *user_entity.LoginUserRequest) (*user_entity.LoginUserResponse, error) {
	user, err := s.Repository.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	err = password_hasher.VerifyPassword(user.Password, payload.Password)
	if err != nil {
		return nil, user_exception.ErrInvalidPassword
	}

	accessToken, err := auth_token_manager.GenerateToken(2*time.Hour, user.Id)
	if err != nil {
		return nil, err
	}

	return &user_entity.LoginUserResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil
}
