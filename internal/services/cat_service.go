package services

import (
	"context"

	"github.com/danzBraham/cats-social/internal/entities/catentity"
	"github.com/danzBraham/cats-social/internal/repositories"
	"github.com/oklog/ulid/v2"
)

type CatService interface {
	CreateCat(ctx context.Context, userId string, payload *catentity.CreateCatRequest) (*catentity.CreateCatResponse, error)
}

type CatServiceImpl struct {
	CatRepository repositories.CatRepository
}

func NewCatService(catRepository repositories.CatRepository) CatService {
	return &CatServiceImpl{CatRepository: catRepository}
}

func (s *CatServiceImpl) CreateCat(ctx context.Context, userId string, payload *catentity.CreateCatRequest) (*catentity.CreateCatResponse, error) {
	cat := &catentity.Cat{
		Id:          ulid.Make().String(),
		Name:        payload.Name,
		Race:        payload.Race,
		Sex:         payload.Sex,
		AgeInMonth:  payload.AgeInMonth,
		Description: payload.Description,
		ImageUrls:   payload.ImageUrls,
		OwnerId:     userId,
	}

	cat, err := s.CatRepository.CreateCat(ctx, cat)
	if err != nil {
		return nil, err
	}

	return &catentity.CreateCatResponse{
		Id:        cat.Id,
		CreatedAt: cat.CreatedAt,
	}, nil
}
