package repositories

import (
	"context"
	"errors"

	"github.com/danzBraham/cats-social/internal/entities/user_entity"
	"github.com/danzBraham/cats-social/internal/errors/user_error"
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
	query := `SELECT 1 FROM users WHERE email = $1`
	var result int
	err := r.DB.QueryRow(ctx, query, email).Scan(&result)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *user_entity.User) error {
	query := `
		INSERT INTO users (id, name, email, password)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.DB.Exec(ctx, query, &user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*user_entity.User, error) {
	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	var user user_entity.User
	err := r.DB.QueryRow(ctx, query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, user_error.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
