package services

import (
	"context"
	"time"

	user_entity "github.com/danzBraham/cats-social/internal/entities/user"
	user_error "github.com/danzBraham/cats-social/internal/errors/user"
	"github.com/danzBraham/cats-social/internal/helpers/bcrypt"
	"github.com/danzBraham/cats-social/internal/helpers/jwt"
	"github.com/danzBraham/cats-social/internal/repositories"
	"github.com/oklog/ulid/v2"
)

type UserService interface {
	RegisterUser(ctx context.Context, payload *user_entity.RegisterUserRequest) (*user_entity.RegisterUserResponse, error)
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

	token, err := jwt.GenerateToken(2*time.Hour, user.Id)
	if err != nil {
		return nil, err
	}

	return &user_entity.RegisterUserResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: token,
	}, nil
}
