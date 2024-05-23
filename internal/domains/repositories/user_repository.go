package repositories

import (
	"context"

	user_entity "github.com/danzbraham/cats-social/internal/domains/entities/users"
)

type UserRepository interface {
	VerifyEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *user_entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*user_entity.User, error)
}
