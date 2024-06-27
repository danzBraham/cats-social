package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/danzBraham/cats-social/internal/entities/matchentity"
	"github.com/danzBraham/cats-social/internal/errors/matcherror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchRepository interface {
	VerifyMatchId(ctx context.Context, matchId string) (bool, error)
	VerifyMatchIdValidity(ctx context.Context, matchId string) (bool, error)
	VerifyRequestIssuer(ctx context.Context, matchId, issuerId string) (bool, error)
	VerifyBothCatsGender(ctx context.Context, matchCatId, userCatId string) (bool, error)
	VerifyBothCatsNotMatched(ctx context.Context, matchCatId, userCatId string) error
	VerifyBothCatsHaveTheSameOwner(ctx context.Context, matchCatId, userCatId string) (bool, error)
	CreateMatch(ctx context.Context, matchCat *matchentity.MatchCat) error
	GetMatches(ctx context.Context, issuerId string) ([]*matchentity.MatchCat, error)
	ApproveMatch(ctx context.Context, matchId string) error
	RejectMatch(ctx context.Context, matchId string) error
}

type MatchRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewMatchRepository(db *pgxpool.Pool) MatchRepository {
	return &MatchRepositoryImpl{DB: db}
}

func (r *MatchRepositoryImpl) VerifyMatchId(ctx context.Context, matchId string) (bool, error) {
	query := `
		SELECT 1
		FROM match_cats
		WHERE id = $1
	`
	var exists int
	err := r.DB.QueryRow(ctx, query, matchId).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MatchRepositoryImpl) VerifyMatchIdValidity(ctx context.Context, matchId string) (bool, error) {
	query := `
		SELECT 1
		FROM match_cats
		WHERE id = $1
			AND status = 'pending'
			AND is_deleted = false
	`
	var exists int
	err := r.DB.QueryRow(ctx, query, matchId).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MatchRepositoryImpl) VerifyRequestIssuer(ctx context.Context, matchId, issuerId string) (bool, error) {
	query := `
		SELECT 1
		FROM match_cats
		WHERE id = $1
			AND issued_by = $2
	`
	var exists int
	err := r.DB.QueryRow(ctx, query, matchId, issuerId).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MatchRepositoryImpl) VerifyBothCatsGender(ctx context.Context, matchCatId, userCatId string) (bool, error) {
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

func (r *MatchRepositoryImpl) VerifyBothCatsNotMatched(ctx context.Context, matchCatId, userCatId string) error {
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
		return matcherror.ErrBothCatsHaveAlreadyMatched
	}
	return nil
}

func (r *MatchRepositoryImpl) VerifyBothCatsHaveTheSameOwner(ctx context.Context, matchCatId, userCatId string) (bool, error) {
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

func (r *MatchRepositoryImpl) CreateMatch(ctx context.Context, matchCat *matchentity.MatchCat) error {
	query := `
		INSERT INTO match_requests (id, match_cat_id, user_cat_id, message)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.DB.Exec(ctx, query,
		&matchCat.Id,
		&matchCat.MatchCatId,
		&matchCat.UserCatId,
		&matchCat.Message,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *MatchRepositoryImpl) GetMatches(ctx context.Context, issuerId string) ([]*matchentity.MatchCat, error) {
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

	matchCats := []*matchentity.MatchCat{}
	for rows.Next() {
		var matchCat matchentity.MatchCat
		var createdAt time.Time
		err := rows.Scan(
			&matchCat.Id,
			&matchCat.MatchCatId,
			&matchCat.UserCatId,
			&matchCat.Message,
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

func (r *MatchRepositoryImpl) ApproveMatch(ctx context.Context, matchId string) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// approve the match cat request
	approveQuery := `
		UPDATE match_cats
		SET status = 'approved'
		WHERE id = $1
		RETURNING match_cat_id, user_cat_id
	`
	var matchCatId, userCatId string
	err = tx.QueryRow(ctx, approveQuery, matchId).Scan(&matchCatId, &userCatId)
	if err != nil {
		return err
	}

	// update has_matched column for both cats
	updateCatsQuery := `
		UPDATE cats
		SET has_matched = true
		WHERE id IN ($1, $2)
	`
	_, err = tx.Exec(ctx, updateCatsQuery, matchCatId, userCatId)
	if err != nil {
		return err
	}

	// remove other match requests for the involved cats
	removeOtherMatchRequestQuery := `
		UPDATE match_cats
		SET is_deleted = true
		WHERE (match_cat_id = $1 OR user_cat_id = $1 OR match_cat_id = $2 OR user_cat_id = $2)
			AND status != 'approved'
	`
	_, err = tx.Exec(ctx, removeOtherMatchRequestQuery, matchCatId, userCatId)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *MatchRepositoryImpl) RejectMatch(ctx context.Context, matchId string) error {
	query := `
		UPDATE match_cats
		SET status = 'rejected'
		WHERE id = $1
	`
	_, err := r.DB.Exec(ctx, query, matchId)
	if err != nil {
		return err
	}
	return nil
}
