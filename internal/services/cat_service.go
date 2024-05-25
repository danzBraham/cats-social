package services

import (
	"context"

	cat_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/cats"
	cat_entity "github.com/danzbraham/cats-social/internal/entities/cat"
	"github.com/danzbraham/cats-social/internal/repositories"
	"github.com/oklog/ulid/v2"
)

type CatService interface {
	AddCat(ctx context.Context, payload *cat_entity.AddCatRequest) (*cat_entity.AddCatResponse, error)
	GetCats(ctx context.Context, params *cat_entity.CatQueryParams) ([]*cat_entity.GetCatReponse, error)
	UpdateCat(ctx context.Context, payload *cat_entity.UpdateCatRequest) error
	DeleteCat(ctx context.Context, payload *cat_entity.DeleteCatRequest) error
}

type CatServiceImpl struct {
	Repository repositories.CatRepository
}

func NewCatService(repository repositories.CatRepository) CatService {
	return &CatServiceImpl{Repository: repository}
}

func (s *CatServiceImpl) AddCat(ctx context.Context, payload *cat_entity.AddCatRequest) (*cat_entity.AddCatResponse, error) {
	cat := &cat_entity.Cat{
		Id:          ulid.Make().String(),
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
		Id:        cat.Id,
		CreatedAt: createdAt,
	}, nil
}

func (s *CatServiceImpl) GetCats(ctx context.Context, params *cat_entity.CatQueryParams) ([]*cat_entity.GetCatReponse, error) {
	return s.Repository.GetCats(ctx, params)
}

func (s *CatServiceImpl) UpdateCat(ctx context.Context, payload *cat_entity.UpdateCatRequest) error {
	isIdExists, err := s.Repository.VerifyId(ctx, payload.Id)
	if err != nil {
		return err
	}
	if !isIdExists {
		return cat_exception.ErrCatIdIsNotFound
	}

	err = s.Repository.UpdateCat(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *CatServiceImpl) DeleteCat(ctx context.Context, payload *cat_entity.DeleteCatRequest) error {
	isIdExists, err := s.Repository.VerifyId(ctx, payload.Id)
	if err != nil {
		return err
	}
	if !isIdExists {
		return cat_exception.ErrCatIdIsNotFound
	}

	err = s.Repository.DeleteCat(ctx, payload.Id)
	if err != nil {
		return err
	}

	return nil
}
