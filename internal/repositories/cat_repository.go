package repositories

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	cat_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/cats"
	cat_entity "github.com/danzbraham/cats-social/internal/entities/cat"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatRepository interface {
	VerifyId(ctx context.Context, id string) (bool, error)
	CreateCat(ctx context.Context, cat *cat_entity.Cat) (createdAt string, err error)
	GetCats(ctx context.Context, userId string, params *cat_entity.CatQueryParams) ([]*cat_entity.GetCatReponse, error)
	GetCatById(ctx context.Context, id string) (*cat_entity.Cat, error)
	GetCatByOwnerId(ctx context.Context, id string) (*cat_entity.Cat, error)
	UpdateCatById(ctx context.Context, id string, cat *cat_entity.UpdateCatRequest) error
	DeleteCatById(ctx context.Context, id string) error
}

type CatRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewCatRepository(db *pgxpool.Pool) CatRepository {
	return &CatRepositoryImpl{DB: db}
}

func (r *CatRepositoryImpl) VerifyId(ctx context.Context, id string) (bool, error) {
	var isIdExists int
	query := `SELECT 1 FROM cats WHERE id = $1 AND is_deleted = false`
	err := r.DB.QueryRow(ctx, query, id).Scan(&isIdExists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *CatRepositoryImpl) CreateCat(ctx context.Context, cat *cat_entity.Cat) (createdAt string, err error) {
	var timeCreated time.Time
	query := `INSERT INTO cats (id, name, race, sex, age_in_month, description, image_urls, owner_id)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
							RETURNING created_at`
	err = r.DB.QueryRow(ctx, query,
		&cat.Id,
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

func (r *CatRepositoryImpl) GetCats(ctx context.Context, userId string, params *cat_entity.CatQueryParams) ([]*cat_entity.GetCatReponse, error) {
	query := `SELECT id, name, race, sex, age_in_month, image_urls, description, has_matched, created_at
						FROM cats 
						WHERE is_deleted = false`
	args := []interface{}{}
	argId := 1

	if params.Id != "" {
		query += ` AND id = $` + strconv.Itoa(argId)
		args = append(args, params.Id)
		argId++
	}

	switch params.Race {
	case "Persian":
		query += ` AND race = 'Persian'`
	case "Maine Coon":
		query += ` AND race = 'Maine Coon'`
	case "Siamese":
		query += ` AND race = 'Siamese'`
	case "Ragdoll":
		query += ` AND race = 'Ragdoll'`
	case "Bengal":
		query += ` AND race = 'Bengal'`
	case "Sphynx":
		query += ` AND race = 'Sphynx'`
	case "British Shorthair":
		query += ` AND race = 'British Shorthair'`
	case "Abyssinian":
		query += ` AND race = 'Abyssinian'`
	case "Scottish Fold":
		query += ` AND race = 'Scottish Fold'`
	case "Birman":
		query += ` AND race = 'Birman'`
	}

	switch params.Sex {
	case "male":
		query += ` AND sex = 'male'`
	case "female":
		query += ` AND sex = 'female'`
	}

	if params.HasMatched != "" {
		hasMatched, err := strconv.ParseBool(params.HasMatched)
		if err != nil {
			return nil, err
		}
		query += ` AND has_matched = $` + strconv.Itoa(argId)
		args = append(args, hasMatched)
		argId++
	}

	if params.AgeInMonth != "" {
		var ageCondition string
		switch {
		case strings.HasPrefix(params.AgeInMonth, ">"):
			ageCondition = `>`
		case strings.HasPrefix(params.AgeInMonth, "<"):
			ageCondition = `<`
		default:
			ageCondition = `=`
		}

		ageValue := strings.TrimLeft(params.AgeInMonth, ">=<")
		age, err := strconv.Atoi(ageValue)
		if err != nil {
			return nil, err
		}

		query += ` AND age_in_month ` + ageCondition + ` $` + strconv.Itoa(argId)
		args = append(args, age)
		argId++
	}

	if params.Owned != "" {
		query += ` AND owner_id = $` + strconv.Itoa(argId)
		args = append(args, userId)
		argId++
	}

	if params.Search != "" {
		query += ` AND name ILIKE $` + strconv.Itoa(argId)
		args = append(args, "%"+params.Search+"%")
		argId++
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argId) + ` OFFSET $` + strconv.Itoa(argId+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cats := []*cat_entity.GetCatReponse{}
	for rows.Next() {
		var cat cat_entity.GetCatReponse
		var timeCreated time.Time
		err := rows.Scan(
			&cat.Id,
			&cat.Name,
			&cat.Race,
			&cat.Sex,
			&cat.AgeInMonth,
			&cat.ImageUrls,
			&cat.Description,
			&cat.HasMatched,
			&timeCreated,
		)
		if err != nil {
			return nil, err
		}
		cat.CreatedAt = timeCreated.Format(time.RFC3339)
		cats = append(cats, &cat)
	}

	return cats, nil
}

func (r *CatRepositoryImpl) GetCatById(ctx context.Context, id string) (*cat_entity.Cat, error) {
	var cat cat_entity.Cat
	query := `SELECT id, name, race, sex, age_in_month, description, image_urls, has_matched, owner_id
						FROM cats
						WHERE id = $1 AND is_deleted = false`
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&cat.Id,
		&cat.Name,
		&cat.Race,
		&cat.Sex,
		&cat.AgeInMonth,
		&cat.Description,
		&cat.ImageUrls,
		&cat.HasMatched,
		&cat.OwnerId,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, cat_exception.ErrCatNotFound
	}
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CatRepositoryImpl) GetCatByOwnerId(ctx context.Context, id string) (*cat_entity.Cat, error) {
	var cat cat_entity.Cat
	query := `SELECT id, name, race, sex, age_in_month, description, image_urls, has_matched, owner_id
						FROM cats
						WHERE owner_id = $1 AND is_deleted = false`
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&cat.Id,
		&cat.Name,
		&cat.Race,
		&cat.Sex,
		&cat.AgeInMonth,
		&cat.Description,
		&cat.ImageUrls,
		&cat.HasMatched,
		&cat.OwnerId,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, cat_exception.ErrCatNotFound
	}
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CatRepositoryImpl) UpdateCatById(ctx context.Context, id string, cat *cat_entity.UpdateCatRequest) error {
	query := `UPDATE cats 
						SET name = $1, race = $2, sex = $3, age_in_month = $4, description = $5, image_urls = $6, updated_at = NOW()
						WHERE id = $7`
	_, err := r.DB.Exec(ctx, query, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, &cat.Description, &cat.ImageUrls, &id)
	if err != nil {
		return err
	}
	return nil
}

func (r *CatRepositoryImpl) DeleteCatById(ctx context.Context, id string) error {
	query := `UPDATE cats SET is_deleted = true WHERE id = $1`
	_, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
