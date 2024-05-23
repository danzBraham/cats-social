package usecases

import (
	"context"
	"time"

	cat_entity "github.com/danzbraham/cats-social/internal/domains/entities/cats"
	"github.com/danzbraham/cats-social/internal/domains/repositories"
	"github.com/oklog/ulid/v2"
)

type CatUsecase interface {
	AddCat(ctx context.Context, payload *cat_entity.AddCatRequest) (*cat_entity.AddCatResponse, error)
}

type CatUsecaseImpl struct {
	Repository repositories.CatRepository
}

func NewCatUsecase(repository repositories.CatRepository) CatUsecase {
	return &CatUsecaseImpl{Repository: repository}
}

func (uc *CatUsecaseImpl) AddCat(ctx context.Context, payload *cat_entity.AddCatRequest) (*cat_entity.AddCatResponse, error) {
	now := time.Now().Format(time.RFC3339)

	cat := &cat_entity.Cat{
		ID:          ulid.Make().String(),
		Name:        payload.Name,
		Race:        payload.Race,
		Sex:         payload.Sex,
		AgeInMonth:  payload.AgeInMonth,
		Description: payload.Description,
		ImageUrls:   payload.ImageUrls,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.Repository.CreateCat(ctx, cat); err != nil {
		return nil, err
	}

	return &cat_entity.AddCatResponse{
		ID:        cat.ID,
		CreatedAt: cat.CreatedAt,
	}, nil
}
