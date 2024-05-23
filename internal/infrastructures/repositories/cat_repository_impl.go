package repositories_impl

import (
	"context"

	cat_entity "github.com/danzbraham/cats-social/internal/domains/entities/cats"
	"github.com/danzbraham/cats-social/internal/domains/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewCatRepositoryImpl(db *pgxpool.Pool) repositories.CatRepository {
	return &CatRepositoryImpl{DB: db}
}

func (r *CatRepositoryImpl) CreateCat(ctx context.Context, cat *cat_entity.Cat) error {
	query := `INSERT INTO cats (id, name, race, sex, age_in_month, description, image_urls)
							VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.DB.Exec(ctx, query, &cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, &cat.Description, &cat.ImageUrls)
	if err != nil {
		return err
	}
	return nil
}
