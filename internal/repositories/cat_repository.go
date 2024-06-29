package repositories

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/danzBraham/cats-social/internal/entities/catentity"
	"github.com/danzBraham/cats-social/internal/errors/caterror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatRepository interface {
	IsCatIdExists(ctx context.Context, catId string) (bool, error)
	IsCatOwner(ctx context.Context, catId, ownerId string) (bool, error)
	CreateCat(ctx context.Context, cat *catentity.Cat) (string, error)
	GetCats(ctx context.Context, ownerId string, params *catentity.CatQueryParams) ([]*catentity.GetCatResponse, error)
	GetCatById(ctx context.Context, catId string) (*catentity.Cat, error)
	UpdateCatById(ctx context.Context, catId string, cat *catentity.Cat) error
	DeleteCatById(ctx context.Context, catId string) error
}

type CatRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewCatRepository(db *pgxpool.Pool) CatRepository {
	return &CatRepositoryImpl{DB: db}
}

func (r *CatRepositoryImpl) IsCatIdExists(ctx context.Context, catId string) (bool, error) {
	query := `
		SELECT
			1
		FROM
			cats
		WHERE
			id = $1
			AND is_deleted = false
	`
	var exists int
	err := r.DB.QueryRow(ctx, query, catId).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *CatRepositoryImpl) IsCatOwner(ctx context.Context, catId, ownerId string) (bool, error) {
	query := `
		SELECT
			1
		FROM
			cats
		WHERE
			id = $1
			AND owner_id = $2
			AND is_deleted = false
	`
	var exists int
	err := r.DB.QueryRow(ctx, query, catId, ownerId).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *CatRepositoryImpl) CreateCat(ctx context.Context, cat *catentity.Cat) (string, error) {
	query := `
		INSERT INTO
			cats (id, name, race, sex, age_in_month, description, image_urls, owner_id)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			created_at
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
		return "", err
	}
	return createdAt.Format(time.RFC3339), nil
}

func (r *CatRepositoryImpl) GetCats(ctx context.Context, ownerId string, params *catentity.CatQueryParams) ([]*catentity.GetCatResponse, error) {
	query := `
		SELECT
			id,
			name,
			race,
			sex,
			age_in_month,
			description,
			image_urls,
			has_matched,
			created_at
		FROM
			cats
		WHERE
			is_deleted = false
	`
	args := []interface{}{}
	argId := 1

	if params.Id != "" {
		query += ` AND id = $` + strconv.Itoa(argId)
		args = append(args, params.Id)
		argId++
	}

	validRace := map[catentity.Race]bool{
		catentity.Persian:          true,
		catentity.MaineCoon:        true,
		catentity.Siamese:          true,
		catentity.Ragdoll:          true,
		catentity.Bengal:           true,
		catentity.Sphynx:           true,
		catentity.BritishShorthair: true,
		catentity.Abyssinian:       true,
		catentity.ScottishFold:     true,
		catentity.Birman:           true,
	}

	if params.Race != "" {
		if !validRace[params.Race] {
			return []*catentity.GetCatResponse{}, nil
		}
		query += ` AND race = $` + strconv.Itoa(argId)
		args = append(args, params.Race)
		argId++
	}

	validSex := map[catentity.Sex]bool{
		catentity.Male:   true,
		catentity.Female: true,
	}

	if params.Sex != "" {
		if !validSex[params.Sex] {
			return []*catentity.GetCatResponse{}, nil
		}
		query += ` AND sex = $` + strconv.Itoa(argId)
		args = append(args, params.Sex)
		argId++
	}

	if params.HasMatched {
		query += ` AND has_matched = $` + strconv.Itoa(argId)
		args = append(args, params.HasMatched)
		argId++
	}

	if params.AgeInMonth != "" {
		var ageCondition string
		switch {
		case strings.HasPrefix(params.AgeInMonth, "<"):
			ageCondition = "<"
		case strings.HasPrefix(params.AgeInMonth, ">"):
			ageCondition = ">"
		default:
			ageCondition = "="
		}

		ageValue := strings.TrimLeft(params.AgeInMonth, "<=>")
		age, err := strconv.Atoi(ageValue)
		if err != nil {
			return nil, err
		}

		query += ` AND age_in_month ` + ageCondition + ` $` + strconv.Itoa(argId)
		args = append(args, age)
		argId++
	}

	if params.Owned {
		query += ` AND owner_id = $` + strconv.Itoa(argId)
		args = append(args, ownerId)
		argId++
	}

	if params.Search != "" {
		query += ` AND name ILIKE $` + strconv.Itoa(argId)
		args = append(args, `%`+params.Search+`%`)
		argId++
	}

	query += ` ORDER BY updated_at DESC LIMIT $` + strconv.Itoa(argId) + ` OFFSET $` + strconv.Itoa(argId+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cats := make([]*catentity.GetCatResponse, 0, params.Limit)
	for rows.Next() {
		var cat catentity.GetCatResponse
		var createdAt time.Time
		err := rows.Scan(
			&cat.Id,
			&cat.Name,
			&cat.Race,
			&cat.Sex,
			&cat.AgeInMonth,
			&cat.Description,
			&cat.ImageUrls,
			&cat.HasMatched,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		cat.CreatedAt = createdAt.Format(time.RFC3339)
		cats = append(cats, &cat)
	}

	return cats, nil
}

func (r *CatRepositoryImpl) GetCatById(ctx context.Context, catId string) (*catentity.Cat, error) {
	query := `
		SELECT
			id,
			name,
			race,
			sex,
			age_in_month,
			description,
			image_urls,
			has_matched,
			owner_id,
			created_at
		FROM
			cats
		WHERE
			id = $1
			AND is_deleted = false
	`
	var cat catentity.Cat
	var createdAt time.Time
	err := r.DB.QueryRow(ctx, query, catId).Scan(
		&cat.Id,
		&cat.Name,
		&cat.Race,
		&cat.Sex,
		&cat.AgeInMonth,
		&cat.Description,
		&cat.ImageUrls,
		&cat.HasMatched,
		&cat.OwnerId,
		&createdAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, caterror.ErrCatNotFound
	}
	if err != nil {
		return nil, err
	}
	cat.CreatedAt = createdAt.Format(time.RFC3339)
	return &cat, nil
}

func (r *CatRepositoryImpl) UpdateCatById(ctx context.Context, catId string, cat *catentity.Cat) error {
	query := `
		UPDATE 
			cats
		SET 
			name = $1,
			race = $2, 
			sex = $3, 
			age_in_month = $4, 
			description = $5, 
			image_urls = $6,
			updated_at = NOW()
		WHERE
			id = $7
			AND is_deleted = false
	`
	_, err := r.DB.Exec(ctx, query,
		&cat.Name,
		&cat.Race,
		&cat.Sex,
		&cat.AgeInMonth,
		&cat.Description,
		&cat.ImageUrls,
		catId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *CatRepositoryImpl) DeleteCatById(ctx context.Context, catId string) error {
	query := `
		UPDATE
			cats
		SET
			is_deleted = true
		WHERE
			id = $1
	`
	_, err := r.DB.Exec(ctx, query, catId)
	if err != nil {
		return err
	}
	return nil
}
