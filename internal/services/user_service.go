package services

import (
	"context"
	"time"

	"github.com/danzBraham/cats-social/internal/entities/user_entity"
	"github.com/danzBraham/cats-social/internal/errors/user_error"
	"github.com/danzBraham/cats-social/internal/helpers/bcrypt"
	"github.com/danzBraham/cats-social/internal/helpers/jwt"
	"github.com/danzBraham/cats-social/internal/repositories"
	"github.com/oklog/ulid/v2"
)

type UserService interface {
	RegisterUser(ctx context.Context, payload *user_entity.RegisterUserRequest) (*user_entity.RegisterUserResponse, error)
	LoginUser(ctx context.Context, payload *user_entity.LoginUserRequest) (*user_entity.LoginUserResponse, error)
}

type UserServiceImpl struct {
	UserRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &UserServiceImpl{UserRepository: userRepository}
}

func (s *UserServiceImpl) RegisterUser(ctx context.Context, payload *user_entity.RegisterUserRequest) (*user_entity.RegisterUserResponse, error) {
	isEmailExists, err := s.UserRepository.VerifyEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}
	if isEmailExists {
		return nil, user_error.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	user := &user_entity.User{
		Id:       ulid.Make().String(),
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
	}

	err = s.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(8*time.Hour, user.Id)
	if err != nil {
		return nil, err
	}

	return &user_entity.RegisterUserResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: token,
	}, nil
}

func (s *UserServiceImpl) LoginUser(ctx context.Context, payload *user_entity.LoginUserRequest) (*user_entity.LoginUserResponse, error) {
	user, err := s.UserRepository.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.VerifyPassword(user.Password, payload.Password)
	if err != nil {
		return nil, user_error.ErrInvalidPassword
	}

	token, err := jwt.GenerateToken(8*time.Hour, user.Id)
	if err != nil {
		return nil, err
	}

	return &user_entity.LoginUserResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: token,
	}, nil
}
