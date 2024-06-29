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
	UpdateCatById(ctx context.Context, userId, catId string, payload *catentity.UpdateCatRequest) error
	DeleteCatById(ctx context.Context, userId, catId string) error
}

type CatServiceImpl struct {
	CatRepository   repositories.CatRepository
	MatchRepository repositories.MatchRepository
}

func NewCatService(catRepository repositories.CatRepository, matchRepository repositories.MatchRepository) CatService {
	return &CatServiceImpl{
		CatRepository:   catRepository,
		MatchRepository: matchRepository,
	}
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

	createdAt, err := s.CatRepository.CreateCat(ctx, cat)
	if err != nil {
		return nil, err
	}

	return &catentity.CreateCatResponse{
		Id:        cat.Id,
		CreatedAt: createdAt,
	}, nil
}

func (s *CatServiceImpl) GetCats(ctx context.Context, userId string, params *catentity.CatQueryParams) ([]*catentity.GetCatResponse, error) {
	return s.CatRepository.GetCats(ctx, userId, params)
}

func (s *CatServiceImpl) UpdateCatById(ctx context.Context, userId, catId string, payload *catentity.UpdateCatRequest) error {
	IsCatIdExists, err := s.CatRepository.IsCatIdExists(ctx, catId)
	if err != nil {
		return err
	}
	if !IsCatIdExists {
		return caterror.ErrCatIdNotFound
	}

	isCatOwner, err := s.CatRepository.IsCatOwner(ctx, catId, userId)
	if err != nil {
		return err
	}
	if !isCatOwner {
		return caterror.ErrNotCatOwner
	}

	isMatchRequestExists, err := s.MatchRepository.IsMatchRequestExists(ctx, catId, catId)
	if err != nil {
		return err
	}
	if isMatchRequestExists {
		return caterror.ErrSexIsEdited
	}

	cat := &catentity.Cat{
		Name:        payload.Name,
		Race:        payload.Race,
		Sex:         payload.Sex,
		AgeInMonth:  payload.AgeInMonth,
		Description: payload.Description,
		ImageUrls:   payload.ImageUrls,
	}

	err = s.CatRepository.UpdateCatById(ctx, catId, cat)
	if err != nil {
		return err
	}

	return nil
}

func (s *CatServiceImpl) DeleteCatById(ctx context.Context, userId, catId string) error {
	IsCatIdExists, err := s.CatRepository.IsCatIdExists(ctx, catId)
	if err != nil {
		return err
	}
	if !IsCatIdExists {
		return caterror.ErrCatIdNotFound
	}

	isCatOwner, err := s.CatRepository.IsCatOwner(ctx, catId, userId)
	if err != nil {
		return err
	}
	if !isCatOwner {
		return caterror.ErrNotCatOwner
	}

	err = s.CatRepository.DeleteCatById(ctx, catId)
	if err != nil {
		return err
	}

	return nil
}
