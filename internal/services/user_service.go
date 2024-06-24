package services

import (
	"context"
	"time"

	"github.com/danzBraham/cats-social/internal/entities/userentity"
	"github.com/danzBraham/cats-social/internal/errors/usererror"
	"github.com/danzBraham/cats-social/internal/helpers/bcrypt"
	"github.com/danzBraham/cats-social/internal/helpers/jwt"
	"github.com/danzBraham/cats-social/internal/repositories"
	"github.com/oklog/ulid/v2"
)

type UserService interface {
	RegisterUser(ctx context.Context, payload *userentity.RegisterUserRequest) (*userentity.RegisterUserResponse, error)
	LoginUser(ctx context.Context, payload *userentity.LoginUserRequest) (*userentity.LoginUserResponse, error)
}

type UserServiceImpl struct {
	UserRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &UserServiceImpl{UserRepository: userRepository}
}

func (s *UserServiceImpl) RegisterUser(ctx context.Context, payload *userentity.RegisterUserRequest) (*userentity.RegisterUserResponse, error) {
	isEmailExists, err := s.UserRepository.VerifyEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}
	if isEmailExists {
		return nil, usererror.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	user := &userentity.User{
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

	return &userentity.RegisterUserResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: token,
	}, nil
}

func (s *UserServiceImpl) LoginUser(ctx context.Context, payload *userentity.LoginUserRequest) (*userentity.LoginUserResponse, error) {
	user, err := s.UserRepository.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.VerifyPassword(user.Password, payload.Password)
	if err != nil {
		return nil, usererror.ErrInvalidPassword
	}

	token, err := jwt.GenerateToken(8*time.Hour, user.Id)
	if err != nil {
		return nil, err
	}

	return &userentity.LoginUserResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: token,
	}, nil
}
