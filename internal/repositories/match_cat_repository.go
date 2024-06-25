package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/danzBraham/cats-social/internal/entities/matchcatentity"
	"github.com/danzBraham/cats-social/internal/errors/matchcaterror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchCatRepository interface {
	VerifyBothCatsGender(ctx context.Context, matchCatId, userCatId string) (bool, error)
	VerifyBothCatsNotMatched(ctx context.Context, matchCatId, userCatId string) error
	VerifyBothCatsHaveTheSameOwner(ctx context.Context, matchCatId, userCatId string) (bool, error)
	CreateMatchCat(ctx context.Context, matchCat *matchcatentity.MatchCat) error
	GetMatchCats(ctx context.Context, issuerId string) ([]*matchcatentity.MatchCat, error)
}

type MatchCatRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewMatchCatRepository(db *pgxpool.Pool) MatchCatRepository {
	return &MatchCatRepositoryImpl{DB: db}
}

func (r *MatchCatRepositoryImpl) VerifyBothCatsGender(ctx context.Context, matchCatId, userCatId string) (bool, error) {
	query := `
		SELECT c1.sex = c2.sex
		FROM cats c1, cats c2
		WHERE c1.id = $1 AND c2.id = $2
	`
	var result bool
	err := r.DB.QueryRow(ctx, query, matchCatId, userCatId).Scan(&result)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result, nil
}

func (r *MatchCatRepositoryImpl) VerifyBothCatsNotMatched(ctx context.Context, matchCatId, userCatId string) error {
	query := `
		SELECT c1.has_matched, c2.has_matched
		FROM cats c1, cats c2
		WHERE c1.id = $1 AND c2.id = $2
	`
	var hasMatched1, hasMatched2 bool
	err := r.DB.QueryRow(ctx, query, matchCatId, userCatId).Scan(&hasMatched1, &hasMatched2)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	if hasMatched1 && hasMatched2 {
		return matchcaterror.ErrBothCatsHaveAlreadyMatched
	}
	return nil
}

func (r *MatchCatRepositoryImpl) VerifyBothCatsHaveTheSameOwner(ctx context.Context, matchCatId, userCatId string) (bool, error) {
	query := `
		SELECT c1.owner_id = c2.owner_id
		FROM cats c1, cats c2
		WHERE c1.id = $1 AND c2.id = $2
	`
	var result bool
	err := r.DB.QueryRow(ctx, query, matchCatId, userCatId).Scan(&result)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result, nil
}

func (r *MatchCatRepositoryImpl) CreateMatchCat(ctx context.Context, matchCat *matchcatentity.MatchCat) error {
	query := `
		INSERT INTO match_cats (id, match_cat_id, user_cat_id, message, issued_by)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.DB.Exec(ctx, query,
		&matchCat.Id,
		&matchCat.MatchCatId,
		&matchCat.UserCatId,
		&matchCat.Message,
		&matchCat.IssuedBy,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *MatchCatRepositoryImpl) GetMatchCats(ctx context.Context, issuerId string) ([]*matchcatentity.MatchCat, error) {
	query := `
		SELECT id, 
					match_cat_id,
					user_cat_id,
					message,
					issued_by,
					created_at
		FROM match_cats
		WHERE issued_by = $1
			AND is_deleted = false
	`
	rows, err := r.DB.Query(ctx, query, issuerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matchCats := []*matchcatentity.MatchCat{}
	for rows.Next() {
		var matchCat matchcatentity.MatchCat
		var createdAt time.Time
		err := rows.Scan(
			&matchCat.Id,
			&matchCat.MatchCatId,
			&matchCat.UserCatId,
			&matchCat.Message,
			&matchCat.IssuedBy,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		matchCat.CreatedAt = createdAt.Format(time.RFC3339)
		matchCats = append(matchCats, &matchCat)
	}

	return matchCats, nil
}
