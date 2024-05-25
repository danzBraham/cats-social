package repositories

import (
	"context"
	"time"

	cat_entity "github.com/danzbraham/cats-social/internal/entities/cats"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatRepository interface {
	CreateCat(ctx context.Context, cat *cat_entity.Cat) (createdAt string, err error)
}

type CatRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewCatRepository(db *pgxpool.Pool) CatRepository {
	return &CatRepositoryImpl{DB: db}
}

func (r *CatRepositoryImpl) CreateCat(ctx context.Context, cat *cat_entity.Cat) (createdAt string, err error) {
	var timeCreated time.Time
	query := `INSERT INTO cats (id, name, race, sex, age_in_month, description, image_urls, owner_id)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
							RETURNING created_at`
	err = r.DB.QueryRow(ctx, query,
		&cat.ID,
		&cat.Name,
		&cat.Race,
		&cat.Sex,
		&cat.AgeInMonth,
		&cat.Description,
		&cat.ImageUrls,
		&cat.OwnerId,
	).Scan(&timeCreated)
	if err != nil {
		return "", err
	}
	createdAt = timeCreated.Format(time.RFC3339)
	return createdAt, nil
}
