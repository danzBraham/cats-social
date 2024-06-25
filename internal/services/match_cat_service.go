package services

import (
	"context"

	"github.com/danzBraham/cats-social/internal/entities/matchcatentity"
	"github.com/danzBraham/cats-social/internal/errors/matchcaterror"
	"github.com/danzBraham/cats-social/internal/repositories"
	"github.com/oklog/ulid/v2"
)

type MatchCatService interface {
	CreateMatchCat(ctx context.Context, userId string, payload *matchcatentity.CreateMatchCatRequest) error
}

type MatchCatServiceImpl struct {
	MatchCatRepository repositories.MatchCatRepository
	CatRepository      repositories.CatRepository
}

func NewMatchCatService(
	matchCatRepository repositories.MatchCatRepository,
	catRepository repositories.CatRepository,
) MatchCatService {
	return &MatchCatServiceImpl{
		MatchCatRepository: matchCatRepository,
		CatRepository:      catRepository,
	}
}

func (s *MatchCatServiceImpl) CreateMatchCat(ctx context.Context, userId string, payload *matchcatentity.CreateMatchCatRequest) error {
	isMatchCatIdExists, err := s.CatRepository.VerifyId(ctx, payload.MatchCatId)
	if err != nil {
		return err
	}
	if !isMatchCatIdExists {
		return matchcaterror.ErrMatchCatIdNotFound
	}

	isUserCatIdExists, err := s.CatRepository.VerifyId(ctx, payload.UserCatId)
	if err != nil {
		return err
	}
	if !isUserCatIdExists {
		return matchcaterror.ErrUserCatIdNotFound
	}

	isUserCatIdBelongToTheUser, err := s.CatRepository.VerifyOwner(ctx, payload.UserCatId, userId)
	if err != nil {
		return err
	}
	if !isUserCatIdBelongToTheUser {
		return matchcaterror.ErrUserCatIdNotBelongToTheUser
	}

	isBothCatsHaveTheSameGender, err := s.MatchCatRepository.VerifyBothCatsGender(ctx, payload.MatchCatId, payload.UserCatId)
	if err != nil {
		return err
	}
	if isBothCatsHaveTheSameGender {
		return matchcaterror.ErrBothCatsHaveTheSameGender
	}

	err = s.MatchCatRepository.VerifyBothCatsNotMatched(ctx, payload.MatchCatId, payload.UserCatId)
	if err != nil {
		return err
	}

	isBothCatsHaveTheSameOwner, err := s.MatchCatRepository.VerifyBothCatsHaveTheSameOwner(ctx, payload.MatchCatId, payload.UserCatId)
	if err != nil {
		return err
	}
	if isBothCatsHaveTheSameOwner {
		return matchcaterror.ErrBothCatsHaveTheSameOwner
	}

	matchCat := &matchcatentity.MatchCat{
		Id:         ulid.Make().String(),
		MatchCatId: payload.MatchCatId,
		UserCatId:  payload.UserCatId,
		Message:    payload.Message,
		IssuedBy:   userId,
	}

	err = s.MatchCatRepository.CreateMatchCat(ctx, matchCat)
	if err != nil {
		return err
	}

	return nil
}
