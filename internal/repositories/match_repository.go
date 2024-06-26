package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/danzBraham/cats-social/internal/entities/matchentity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchRepository interface {
	IsMatchIdExists(ctx context.Context, matchId string) (bool, error)
	IsMatchIdValid(ctx context.Context, matchId string) (bool, error)
	IsMatchIssuer(ctx context.Context, matchId, userId string) (bool, error)
	IsBothCatsHaveSameGender(ctx context.Context, matchCatId, userCatId string) (bool, error)
	IsBothCatsAlreadyMatched(ctx context.Context, matchCatId, userCatId string) (bool, error)
	IsOwnerOfBothCats(ctx context.Context, matchCatId, userCatId string) (bool, error)
	IsMatchRequestExists(ctx context.Context, matchCatId, userCatId string) (bool, error)
	CreateMatch(ctx context.Context, matchCat *matchentity.Match) error
	GetMatches(ctx context.Context, userId string) ([]*matchentity.GetMatchResponse, error)
	ApproveMatch(ctx context.Context, matchId string) error
	RejectMatch(ctx context.Context, matchId string) error
	DeleteMatch(ctx context.Context, matchId string) error
}

type MatchRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewMatchRepository(db *pgxpool.Pool) MatchRepository {
	return &MatchRepositoryImpl{DB: db}
}

func (r *MatchRepositoryImpl) IsMatchIdExists(ctx context.Context, matchId string) (bool, error) {
	query := `
		SELECT
			1
		FROM
			match_requests
		WHERE
			id = $1
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

func (r *MatchRepositoryImpl) IsMatchIdValid(ctx context.Context, matchId string) (bool, error) {
	query := `
		SELECT 
			1
		FROM 
			match_requests
		WHERE 
			id = $1
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

func (r *MatchRepositoryImpl) IsMatchIssuer(ctx context.Context, matchId, userId string) (bool, error) {
	query := `
		SELECT 
			1
		FROM 
			match_requests mr
		JOIN 
			cats c ON mr.user_cat_id = c.id
		WHERE 
			mr.id = $1
			AND c.owner_id = $2
	`
	var exists int
	err := r.DB.QueryRow(ctx, query, matchId, userId).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MatchRepositoryImpl) IsBothCatsHaveSameGender(ctx context.Context, matchCatId, userCatId string) (bool, error) {
	query := `
		SELECT
			c1.sex = c2.sex
		FROM
			cats c1, cats c2
		WHERE
			c1.id = $1
			AND c2.id = $2
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

func (r *MatchRepositoryImpl) IsBothCatsAlreadyMatched(ctx context.Context, matchCatId, userCatId string) (bool, error) {
	query := `
		SELECT
			status = 'approved'
		FROM
			match_requests
		WHERE
			(match_cat_id = $1 OR user_cat_id = $1)
			AND (match_cat_id = $2 OR user_cat_id = $2)
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

func (r *MatchRepositoryImpl) IsOwnerOfBothCats(ctx context.Context, matchCatId, userCatId string) (bool, error) {
	query := `
		SELECT 
			c1.owner_id = c2.owner_id
		FROM
			cats c1, cats c2
		WHERE
			c1.id = $1
			AND c2.id = $2
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

func (r *MatchRepositoryImpl) IsMatchRequestExists(ctx context.Context, matchCatId, userCatId string) (bool, error) {
	query := `
		SELECT
			1
		FROM
			match_requests
		WHERE
			(match_cat_id = $1 OR user_cat_id = $1)
			AND (user_cat_id = $2 OR match_cat_id = $2)
			AND status = 'pending'
			AND is_deleted = false
	`
	var exists int
	err := r.DB.QueryRow(ctx, query, matchCatId, userCatId).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MatchRepositoryImpl) CreateMatch(ctx context.Context, matchCat *matchentity.Match) error {
	query := `
		INSERT INTO
			match_requests (id, match_cat_id, user_cat_id, message)
		VALUES
			($1, $2, $3, $4)
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

func (r *MatchRepositoryImpl) GetMatches(ctx context.Context, userId string) ([]*matchentity.GetMatchResponse, error) {
	query := `
		SELECT
			mr.id,
			u.name AS issuer_name,
			u.email AS issuer_email,
			u.created_at AS issuer_created_at,
			mc.id AS mc_id,
			mc.name AS mc_name,
			mc.race AS mc_race,
			mc.sex AS mc_sex,
			mc.description AS mc_description,
			mc.age_in_month AS mc_age_in_month,
			mc.image_urls AS mc_image_urls,
			mc.has_matched AS mc_has_matched,
			mc.created_at AS mc_created_at,
			uc.id AS uc_id,
			uc.name AS uc_name,
			uc.race AS uc_race,
			uc.sex AS uc_sex,
			uc.description AS uc_description,
			uc.age_in_month AS uc_age_in_month,
			uc.image_urls AS uc_image_urls,
			uc.has_matched AS uc_has_matched,
			uc.created_at AS uc_created_at,
			mr.message,
			mr.created_at
		FROM
			match_requests mr
		JOIN
			cats mc ON mr.match_cat_id = mc.id
		JOIN
			cats uc ON mr.user_cat_id = uc.id
		JOIN
			users u ON uc.owner_id = u.id
		WHERE
			(mc.owner_id = $1 OR uc.owner_id = $1)
			AND mr.is_deleted = false
		ORDER BY
			created_at
	`
	rows, err := r.DB.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := []*matchentity.GetMatchResponse{}
	for rows.Next() {
		var match matchentity.GetMatchResponse
		var issuerCreatedAt, matchCatCreatedAt, userCatCreatedAt, matchCreatedAt time.Time
		err := rows.Scan(
			&match.Id,
			&match.IssuedBy.Name,
			&match.IssuedBy.Email,
			&issuerCreatedAt,
			&match.MatchCatDetail.Id,
			&match.MatchCatDetail.Name,
			&match.MatchCatDetail.Race,
			&match.MatchCatDetail.Sex,
			&match.MatchCatDetail.Description,
			&match.MatchCatDetail.AgeInMonth,
			&match.MatchCatDetail.ImageUrls,
			&match.MatchCatDetail.HasMatched,
			&matchCatCreatedAt,
			&match.UserCatDetail.Id,
			&match.UserCatDetail.Name,
			&match.UserCatDetail.Race,
			&match.UserCatDetail.Sex,
			&match.UserCatDetail.Description,
			&match.UserCatDetail.AgeInMonth,
			&match.UserCatDetail.ImageUrls,
			&match.UserCatDetail.HasMatched,
			&userCatCreatedAt,
			&match.Message,
			&matchCreatedAt,
		)
		if err != nil {
			return nil, err
		}
		match.IssuedBy.CreatedAt = issuerCreatedAt.Format(time.RFC3339)
		match.MatchCatDetail.CreatedAt = matchCatCreatedAt.Format(time.RFC3339)
		match.UserCatDetail.CreatedAt = userCatCreatedAt.Format(time.RFC3339)
		match.CreatedAt = matchCreatedAt.Format(time.RFC3339)
		matches = append(matches, &match)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func (r *MatchRepositoryImpl) ApproveMatch(ctx context.Context, matchId string) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// approve the match request
	approveQuery := `
		UPDATE
			match_requests
		SET
			status = 'approved'
		WHERE 
			id = $1
		RETURNING 
			match_cat_id, user_cat_id
	`
	var matchCatId, userCatId string
	err = tx.QueryRow(ctx, approveQuery, matchId).Scan(&matchCatId, &userCatId)
	if err != nil {
		return err
	}

	// update has_matched column for both cats
	updateCatsQuery := `
		UPDATE
			cats
		SET
			has_matched = true
		WHERE
			id IN ($1, $2)
	`
	_, err = tx.Exec(ctx, updateCatsQuery, matchCatId, userCatId)
	if err != nil {
		return err
	}

	// remove other match requests for the involved cats
	removeOtherMatchRequestQuery := `
		UPDATE
			match_requests
		SET
			is_deleted = true
		WHERE
			(match_cat_id = $1 OR user_cat_id = $1 OR match_cat_id = $2 OR user_cat_id = $2)
			AND status != 'approved'
	`
	_, err = tx.Exec(ctx, removeOtherMatchRequestQuery, matchCatId, userCatId)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *MatchRepositoryImpl) RejectMatch(ctx context.Context, matchId string) error {
	query := `
		UPDATE 
			match_requests
		SET 
			status = 'rejected'
		WHERE 
			id = $1
	`
	_, err := r.DB.Exec(ctx, query, matchId)
	if err != nil {
		return err
	}
	return nil
}

func (r *MatchRepositoryImpl) DeleteMatch(ctx context.Context, matchId string) error {
	query := `
		UPDATE
			match_requests
		SET
			is_deleted = true
		WHERE
			id = $1
	`
	_, err := r.DB.Exec(ctx, query, matchId)
	if err != nil {
		return err
	}
	return nil
}
