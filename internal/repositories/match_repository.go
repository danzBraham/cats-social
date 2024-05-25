package repositories

import (
	"context"
	"time"

	match_entity "github.com/danzbraham/cats-social/internal/entities/match"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchRepository interface {
	CreateMatchCat(ctx context.Context, matchCat *match_entity.MatchCat) error
	GetMatchCatRequests(ctx context.Context) ([]*match_entity.GetMatchCatResponse, error)
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