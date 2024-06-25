package repositories

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/danzBraham/cats-social/internal/entities/catentity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatRepository interface {
	VerifyId(ctx context.Context, id string) (bool, error)
	CreateCat(ctx context.Context, cat *catentity.Cat) (*catentity.Cat, error)
	GetCats(ctx context.Context, userId string, params *catentity.CatQueryParams) ([]*catentity.Cat, error)
	UpdateCatById(ctx context.Context, id string, cat *catentity.Cat) error
}

type CatRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewCatRepository(db *pgxpool.Pool) CatRepository {
	return &CatRepositoryImpl{DB: db}
}

func (r *CatRepositoryImpl) VerifyId(ctx context.Context, id string) (bool, error) {
	query := `SELECT 1 FROM cats WHERE id = $1`
	var result int
	err := r.DB.QueryRow(ctx, query, id).Scan(&result)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
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

func (r *CatRepositoryImpl) GetCats(ctx context.Context, userId string, params *catentity.CatQueryParams) ([]*catentity.Cat, error) {
	query := `
		SELECT id, name, race, sex, age_in_month, description, image_urls, has_matched, created_at
		FROM cats
		WHERE is_deleted = false
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
			return []*catentity.Cat{}, nil
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
			return []*catentity.Cat{}, nil
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
		args = append(args, userId)
		argId++
	}

	if params.Search != "" {
		query += ` AND name ILIKE $` + strconv.Itoa(argId)
		args = append(args, `%`+params.Search+`%`)
		argId++
	}

	query += ` LIMIT $` + strconv.Itoa(argId) + ` OFFSET $` + strconv.Itoa(argId+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cats := make([]*catentity.Cat, 0, params.Limit)
	for rows.Next() {
		var cat catentity.Cat
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

func (r *CatRepositoryImpl) UpdateCatById(ctx context.Context, id string, cat *catentity.Cat) error {
	query := `
		UPDATE cats
		SET name = $1,
				race = $2, 
				sex = $3, 
				age_in_month = $4, 
				description = $5, 
				image_urls = $6,
				updated_at = NOW()
		WHERE id = $7
	`
	_, err := r.DB.Exec(ctx, query,
		&cat.Name,
		&cat.Race,
		&cat.Sex,
		&cat.AgeInMonth,
		&cat.Description,
		&cat.ImageUrls,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
