package services

import (
	"context"

	match_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/match"
	match_entity "github.com/danzbraham/cats-social/internal/entities/match"
	"github.com/danzbraham/cats-social/internal/repositories"
	"github.com/oklog/ulid/v2"
)

type MatchService interface {
	RequestMatchCat(ctx context.Context, payload *match_entity.MatchCatRequest) error
	GetMatchCatRequests(ctx context.Context) ([]*match_entity.GetMatchCatResponse, error)
}

type MatchServiceImpl struct {
	CatRepository   repositories.CatRepository
	MatchRepository repositories.MatchRepository
}

func NewMatchService(catRepository repositories.CatRepository, matchRepository repositories.MatchRepository) MatchService {
	return &MatchServiceImpl{
		CatRepository:   catRepository,
		MatchRepository: matchRepository,
	}
}

func (s *MatchServiceImpl) RequestMatchCat(ctx context.Context, payload *match_entity.MatchCatRequest) error {
	isMatchCatIdExists, err := s.CatRepository.VerifyId(ctx, payload.MatchCatId)
	if err != nil {
		return err
	}
	if !isMatchCatIdExists {
		return match_exception.ErrMatchCatIdIsNotFound
	}

	isUserCatIdExists, err := s.CatRepository.VerifyId(ctx, payload.UserCatId)
	if err != nil {
		return err
	}
	if !isUserCatIdExists {
		return match_exception.ErrUserCatIdIsNotFound
	}

	issuerCat, err := s.CatRepository.GetCatByOwnerId(ctx, payload.Issuer)
	if err != nil {
		return err
	}

	if payload.UserCatId != issuerCat.Id {
		return match_exception.ErrUserCatIdNotBelongTheUser
	}

	matchCat, err := s.CatRepository.GetCatById(ctx, payload.MatchCatId)
	if err != nil {
		return err
	}

	userCat, err := s.CatRepository.GetCatById(ctx, payload.UserCatId)
	if err != nil {
		return err
	}

	if matchCat.Sex == userCat.Sex {
		return match_exception.ErrBothCatHaveSameGender
	}

	if matchCat.HasMatched && userCat.HasMatched {
		return match_exception.ErrBothCatAlreadyMatched
	}

	if matchCat.OwnerId == userCat.OwnerId {
		return match_exception.ErrBothCatsHaveTheSameOwner
	}

	createMatchCat := &match_entity.MatchCat{
		Id:         ulid.Make().String(),
		MatchCatId: payload.MatchCatId,
		UserCatId:  payload.UserCatId,
		Message:    payload.Message,
		Status:     match_entity.Pending,
		IssuedBy:   payload.Issuer,
	}

	err = s.MatchRepository.CreateMatchCat(ctx, createMatchCat)
	if err != nil {
		return err
	}

	return nil
}

func (s *MatchServiceImpl) GetMatchCatRequests(ctx context.Context) ([]*match_entity.GetMatchCatResponse, error) {
	return s.MatchRepository.GetMatchCatRequests(ctx)
}
