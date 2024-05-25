package repositories

import (
	"context"

	match_entity "github.com/danzbraham/cats-social/internal/entities/match"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchRepository interface {
	CreateMatchCat(ctx context.Context, matchCat *match_entity.MatchCat) error
}

type MatchRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewMatchRepository(db *pgxpool.Pool) MatchRepository {
	return &MatchRepositoryImpl{DB: db}
}

func (r *MatchRepositoryImpl) CreateMatchCat(ctx context.Context, matchCat *match_entity.MatchCat) error {
	query := `INSERT INTO match_cats (id, match_cat_id, user_cat_id, message, status, issued_by)
						VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.DB.Exec(ctx, query,
		&matchCat.Id,
		&matchCat.MatchCatId,
		&matchCat.UserCatId,
		&matchCat.Message,
		&matchCat.Status,
		&matchCat.IssuedBy,
	)
	if err != nil {
		return err
	}
	return nil
}
