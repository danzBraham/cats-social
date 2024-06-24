package repositories

import (
	"context"
	"time"

	"github.com/danzBraham/cats-social/internal/entities/catentity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatRepository interface {
	CreateCat(ctx context.Context, cat *catentity.Cat) (*catentity.Cat, error)
}

type CatRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewCatRepository(db *pgxpool.Pool) CatRepository {
	return &CatRepositoryImpl{DB: db}
}

func (r *CatRepositoryImpl) CreateCat(ctx context.Context, cat *catentity.Cat) (*catentity.Cat, error) {
	query := `
		INSERT INTO cats (id, name, race, sex, age_in_month, description, image_urls, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at
	`
	var createdAt time.Time
	err := r.DB.QueryRow(ctx, query,
		&cat.Id,
		&cat.Name,
		&cat.Race,
		&cat.Sex,
		&cat.AgeInMonth,
		&cat.Description,
		&cat.ImageUrls,
		&cat.OwnerId,
	).Scan(&createdAt)
	if err != nil {
		return nil, err
	}
	cat.CreatedAt = createdAt.Format(time.RFC3339)
	return cat, nil
}
