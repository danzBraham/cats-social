package services

import (
	"context"

	cat_entity "github.com/danzbraham/cats-social/internal/entities/cats"
	"github.com/danzbraham/cats-social/internal/repositories"
	"github.com/oklog/ulid/v2"
)

type CatService interface {
	AddCat(ctx context.Context, payload *cat_entity.AddCatRequest) (*cat_entity.AddCatResponse, error)
	GetCats(ctx context.Context, params *cat_entity.CatQueryParams) ([]*cat_entity.GetCatReponse, error)
}

type CatServiceImpl struct {
	Repository repositories.CatRepository
}

func NewCatService(repository repositories.CatRepository) CatService {
	return &CatServiceImpl{Repository: repository}
}

func (s *CatServiceImpl) AddCat(ctx context.Context, payload *cat_entity.AddCatRequest) (*cat_entity.AddCatResponse, error) {
	cat := &cat_entity.Cat{
		ID:          ulid.Make().String(),
		Name:        payload.Name,
		Race:        payload.Race,
		Sex:         payload.Sex,
		AgeInMonth:  payload.AgeInMonth,
		Description: payload.Description,
		ImageUrls:   payload.ImageUrls,
		OwnerId:     payload.OwnerId,
	}

	createdAt, err := s.Repository.CreateCat(ctx, cat)
	if err != nil {
		return nil, err
	}

	return &cat_entity.AddCatResponse{
		ID:        cat.ID,
		CreatedAt: createdAt,
	}, nil
}

func (s *CatServiceImpl) GetCats(ctx context.Context, params *cat_entity.CatQueryParams) ([]*cat_entity.GetCatReponse, error) {
	return s.Repository.GetCats(ctx, params)
}
