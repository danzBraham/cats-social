package repositories

import (
	"context"
	"errors"

	user_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/users"
	user_entity "github.com/danzbraham/cats-social/internal/entities/users"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	VerifyEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *user_entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*user_entity.User, error)
}

type UserRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (r *UserRepositoryImpl) VerifyEmail(ctx context.Context, email string) (bool, error) {
	var isEmailExists int
	query := `SELECT 1 FROM users WHERE email = $1`
	err := r.DB.QueryRow(ctx, query, email).Scan(&isEmailExists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *user_entity.User) error {
	query := `INSERT INTO users (id, email, name, password) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(ctx, query, &user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*user_entity.User, error) {
	user := &user_entity.User{}
	query := `SELECT id, email, name, password FROM users WHERE email = $1`
	err := r.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, user_exception.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
