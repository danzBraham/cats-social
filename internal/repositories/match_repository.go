package repositories

import (
	"context"
	"errors"
	"time"

	match_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/match"
	match_entity "github.com/danzbraham/cats-social/internal/entities/match"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchRepository interface {
	VerifyId(ctx context.Context, matchId string) (bool, bool, error)
	VerifyStatus(ctx context.Context, matchId string) (bool, error)
	VerifyIssuer(ctx context.Context, userId string) (bool, error)
	CreateMatchCatRequest(ctx context.Context, matchCat *match_entity.MatchCat) error
	GetMatchCatRequests(ctx context.Context) ([]*match_entity.GetMatchCatResponse, error)
	GetMatchCatRequestById(ctx context.Context, matchId string) (*match_entity.MatchCat, error)
	ApproveMatchCatRequest(ctx context.Context, matchId string) error
	RejectMatchCatRequest(ctx context.Context, matchId string) error
	DeleteMatchCatRequest(ctx context.Context, matchId string) error
}

type MatchRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewMatchRepository(db *pgxpool.Pool) MatchRepository {
	return &MatchRepositoryImpl{DB: db}
}

func (r *MatchRepositoryImpl) VerifyId(ctx context.Context, matchId string) (bool, bool, error) {
	var isDeleted bool
	query := `SELECT is_deleted FROM match_cats WHERE id = $1`
	err := r.DB.QueryRow(ctx, query, matchId).Scan(&isDeleted)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}
	return true, isDeleted, nil
}

func (r *MatchRepositoryImpl) VerifyStatus(ctx context.Context, matchId string) (bool, error) {
	var isPending int
	query := `SELECT 1 FROM match_cats WHERE id = $1 AND status = 'pending'`
	err := r.DB.QueryRow(ctx, query, matchId).Scan(&isPending)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MatchRepositoryImpl) VerifyIssuer(ctx context.Context, userId string) (bool, error) {
	var isIssuer int
	query := `SELECT 1 FROM match_cats WHERE issued_by = $1`
	err := r.DB.QueryRow(ctx, query, userId).Scan(&isIssuer)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MatchRepositoryImpl) CreateMatchCatRequest(ctx context.Context, matchCat *match_entity.MatchCat) error {
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

func (r *MatchRepositoryImpl) GetMatchCatRequests(ctx context.Context) ([]*match_entity.GetMatchCatResponse, error) {
	query := `SELECT m.id,
									u.name, u.email, u.created_at,
									mc.id, mc.name, mc.race, mc.sex, mc.age_in_month, 
									mc.description, mc.image_urls, mc.has_matched, mc.created_at,
									uc.id, uc.name, uc.race, uc.sex, uc.age_in_month, 
									uc.description, uc.image_urls, uc.has_matched, uc.created_at,
									m.message, m.created_at
						FROM match_cats m
						JOIN cats mc ON m.match_cat_id = mc.id
						JOIN cats uc ON m.user_cat_id = uc.id
						JOIN users u ON m.issued_by = u.id
						WHERE m.is_deleted = false
						ORDER BY m.created_at DESC`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	matchCats := []*match_entity.GetMatchCatResponse{}
	for rows.Next() {
		var matchCat match_entity.GetMatchCatResponse
		var issuerCreatedAt, matchCatCreatedAt, userCatCreatedAt, createdAt time.Time
		err := rows.Scan(
			&matchCat.Id,
			&matchCat.IssuedBy.Name, &matchCat.IssuedBy.Email, &issuerCreatedAt,
			&matchCat.MatchCatDetail.Id, &matchCat.MatchCatDetail.Name, &matchCat.MatchCatDetail.Race, &matchCat.MatchCatDetail.Sex,
			&matchCat.MatchCatDetail.AgeInMonth, &matchCat.MatchCatDetail.Description, &matchCat.MatchCatDetail.ImageUrls,
			&matchCat.MatchCatDetail.HasMatched, &matchCatCreatedAt,
			&matchCat.UserCatDetail.Id, &matchCat.UserCatDetail.Name, &matchCat.UserCatDetail.Race, &matchCat.UserCatDetail.Sex,
			&matchCat.UserCatDetail.AgeInMonth, &matchCat.UserCatDetail.Description, &matchCat.UserCatDetail.ImageUrls,
			&matchCat.UserCatDetail.HasMatched, &userCatCreatedAt,
			&matchCat.Message, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		matchCat.IssuedBy.CreatedAt = issuerCreatedAt.Format(time.RFC3339)
		matchCat.MatchCatDetail.CreatedAt = matchCatCreatedAt.Format(time.RFC3339)
		matchCat.UserCatDetail.CreatedAt = userCatCreatedAt.Format(time.RFC3339)
		matchCat.CreatedAt = createdAt.Format(time.RFC3339)

		matchCats = append(matchCats, &matchCat)
	}

	return matchCats, nil
}

func (r *MatchRepositoryImpl) GetMatchCatRequestById(ctx context.Context, id string) (*match_entity.MatchCat, error) {
	var matchCat match_entity.MatchCat
	var createdAt time.Time
	query := `SELECT id, match_cat_id, user_cat_id, message, status, issued_by, created_at
						FROM match_cats
						WHERE id = $1 AND is_deleted = false`
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&matchCat.Id,
		&matchCat.MatchCatId,
		&matchCat.UserCatId,
		&matchCat.Message,
		&matchCat.Status,
		&matchCat.IssuedBy,
		&createdAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, match_exception.ErrMatchCatIsNotFound
	}
	if err != nil {
		return nil, err
	}
	matchCat.CreatedAt = createdAt.Format(time.RFC3339)
	return &matchCat, nil
}

func (r *MatchRepositoryImpl) ApproveMatchCatRequest(ctx context.Context, matchId string) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Approve the match
	approveMatchQuery := `UPDATE match_cats
												SET status = 'approved', updated_at = NOW()
												WHERE id = $1 AND status != 'approved'`
	_, err = tx.Exec(ctx, approveMatchQuery, matchId)
	if err != nil {
		return err
	}

	// Get match details
	var matchCatId, userCatId string
	getMatchDetailQuery := `SELECT match_cat_id, user_cat_id
													FROM match_cats
													WHERE id = $1`
	err = tx.QueryRow(ctx, getMatchDetailQuery, matchId).Scan(&matchCatId, &userCatId)
	if err != nil {
		return err
	}

	// Remove other match requests involving the same pair of cats
	removeOtherMatchesQuery := `UPDATE match_cats
															SET is_deleted = true, updated_at = NOW()
															WHERE id != $1
																AND status != 'approved'
																AND (
																	(match_cat_id = $2 AND user_cat_id = $3)
																	OR
																	(match_cat_id = $3 AND user_cat_id = $2)
																)`
	_, err = tx.Exec(ctx, removeOtherMatchesQuery, matchId, matchCatId, userCatId)
	if err != nil {
		return err
	}

	// Update has_matched field in the cats table for both cats
	updateHasMatchedQuery := `UPDATE cats
														SET has_matched = true
														WHERE id = $1 OR id = $2`
	_, err = tx.Exec(ctx, updateHasMatchedQuery, matchCatId, userCatId)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *MatchRepositoryImpl) RejectMatchCatRequest(ctx context.Context, matchId string) error {
	query := `UPDATE match_cats
						SET status = 'rejected', updated_at = NOW()
						WHERE id = $1 AND status != 'rejected'`
	_, err := r.DB.Exec(ctx, query, matchId)
	if err != nil {
		return err
	}
	return nil
}

func (r *MatchRepositoryImpl) DeleteMatchCatRequest(ctx context.Context, matchId string) error {
	query := `UPDATE match_cats
						SET is_deleted = true, updated_at = NOW()
						WHERE id = $1`
	_, err := r.DB.Exec(ctx, query, matchId)
	if err != nil {
		return err
	}
	return nil
}
