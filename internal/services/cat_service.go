package services

import (
	"context"

	"github.com/danzBraham/cats-social/internal/entities/catentity"
	"github.com/danzBraham/cats-social/internal/errors/caterror"
	"github.com/danzBraham/cats-social/internal/repositories"
	"github.com/oklog/ulid/v2"
)

type CatService interface {
	CreateCat(ctx context.Context, userId string, payload *catentity.CreateCatRequest) (*catentity.CreateCatResponse, error)
	GetCats(ctx context.Context, userId string, params *catentity.CatQueryParams) ([]*catentity.GetCatResponse, error)
	UpdateCatById(ctx context.Context, id string, payload *catentity.UpdateCatRequest) error
	DeleteCatById(ctx context.Context, id string) error
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

func (s *CatServiceImpl) GetCats(ctx context.Context, userId string, params *catentity.CatQueryParams) ([]*catentity.GetCatResponse, error) {
	cats, err := s.CatRepository.GetCats(ctx, userId, params)
	if err != nil {
		return nil, err
	}

	catsResponse := make([]*catentity.GetCatResponse, 0, len(cats))
	for _, cat := range cats {
		catsResponse = append(catsResponse, &catentity.GetCatResponse{
			Id:          cat.Id,
			Name:        cat.Name,
			Race:        cat.Race,
			Sex:         cat.Sex,
			AgeInMonth:  cat.AgeInMonth,
			Description: cat.Description,
			ImageUrls:   cat.ImageUrls,
			HasMatched:  cat.HasMatched,
			CreatedAt:   cat.CreatedAt,
		})
	}

	return catsResponse, nil
}

func (s *CatServiceImpl) UpdateCatById(ctx context.Context, id string, payload *catentity.UpdateCatRequest) error {
	isIdExists, err := s.CatRepository.VerifyId(ctx, id)
	if err != nil {
		return err
	}
	if !isIdExists {
		return caterror.ErrIdNotFound
	}

	cat := &catentity.Cat{
		Name:        payload.Name,
		Race:        payload.Race,
		Sex:         payload.Sex,
		AgeInMonth:  payload.AgeInMonth,
		Description: payload.Description,
		ImageUrls:   payload.ImageUrls,
	}

	err = s.CatRepository.UpdateCatById(ctx, id, cat)
	if err != nil {
		return err
	}

	return nil
}

func (s *CatServiceImpl) DeleteCatById(ctx context.Context, id string) error {
	isIdExists, err := s.CatRepository.VerifyId(ctx, id)
	if err != nil {
		return err
	}
	if !isIdExists {
		return caterror.ErrIdNotFound
	}

	err = s.CatRepository.DeleteCatById(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
