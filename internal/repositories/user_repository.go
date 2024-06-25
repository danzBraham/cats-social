package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/danzBraham/cats-social/internal/entities/userentity"
	"github.com/danzBraham/cats-social/internal/errors/usererror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	VerifyEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *userentity.User) error
	GetUserByEmail(ctx context.Context, email string) (*userentity.User, error)
	GetUserById(ctx context.Context, id string) (*userentity.User, error)
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

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *userentity.User) error {
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

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*userentity.User, error) {
	query := `
		SELECT id,
					name,
					email,
					password,
					created_at
		FROM users 
		WHERE email = $1
	`
	var user userentity.User
	var createdAt time.Time
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&createdAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, usererror.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserById(ctx context.Context, email string) (*userentity.User, error) {
	query := `
		SELECT id,
					name,
					email,
					password,
					created_at
		FROM users 
		WHERE id = $1
	`
	var user userentity.User
	var createdAt time.Time
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&createdAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, usererror.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
	return &user, nil
}
