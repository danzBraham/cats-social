package repositories

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	cat_entity "github.com/danzbraham/cats-social/internal/entities/cats"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatRepository interface {
	CreateCat(ctx context.Context, cat *cat_entity.Cat) (createdAt string, err error)
	GetCats(ctx context.Context, params *cat_entity.CatQueryParams) ([]*cat_entity.GetCatReponse, error)
}

type CatRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewCatRepository(db *pgxpool.Pool) CatRepository {
	return &CatRepositoryImpl{DB: db}
}

func (r *CatRepositoryImpl) CreateCat(ctx context.Context, cat *cat_entity.Cat) (createdAt string, err error) {
	var timeCreated time.Time
	query := `INSERT INTO cats (id, name, race, sex, age_in_month, description, image_urls, owner_id)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
							RETURNING created_at`
	err = r.DB.QueryRow(ctx, query,
		&cat.ID,
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

func (r *CatRepositoryImpl) GetCats(ctx context.Context, params *cat_entity.CatQueryParams) ([]*cat_entity.GetCatReponse, error) {
	query := `SELECT id, name, race, sex, age_in_month, image_urls, description, has_matched, created_at
							FROM cats 
							WHERE is_deleted = false`
	args := []interface{}{}
	argID := 1

	if params.ID != "" {
		query += ` AND id = $` + strconv.Itoa(argID)
		args = append(args, params.ID)
		argID++
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
		query += ` AND has_matched = $` + strconv.Itoa(argID)
		args = append(args, hasMatched)
		argID++
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

		query += ` AND age_in_month ` + ageCondition + ` $` + strconv.Itoa(argID)
		args = append(args, age)
		argID++
	}

	if params.Owned != "" {
		owned, err := strconv.ParseBool(params.Owned)
		if err != nil {
			return nil, err
		}
		query += ` AND owned = $` + strconv.Itoa(argID)
		args = append(args, owned)
		argID++
	}

	if params.Search != "" {
		query += ` AND name ILIKE $` + strconv.Itoa(argID)
		args = append(args, "%"+params.Search+"%")
		argID++
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argID) + ` OFFSET $` + strconv.Itoa(argID+1)
	args = append(args, params.Limit, params.Offset)
	fmt.Println(query)

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
			&cat.ID,
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
