package repositories_impl

import (
	"context"
	"errors"

	user_entity "github.com/danzbraham/cats-social/internal/domains/entities/users"
	"github.com/danzbraham/cats-social/internal/domains/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryPostgres struct {
	DB *pgxpool.Pool
}

func NewUserRepositoryPostgres(db *pgxpool.Pool) repositories.UserRepository {
	return &UserRepositoryPostgres{DB: db}
}

func (r *UserRepositoryPostgres) VerifyEmail(ctx context.Context, email string) (bool, error) {
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

func (r *UserRepositoryPostgres) CreateUser(ctx context.Context, user *user_entity.User) error {
	query := `INSERT INTO users (id, email, name, password) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(ctx, query, &user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return err
	}
	return nil
}
