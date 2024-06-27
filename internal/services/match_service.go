package services

import (
	"context"

	"github.com/danzBraham/cats-social/internal/entities/matchentity"
	"github.com/danzBraham/cats-social/internal/errors/matcherror"
	"github.com/danzBraham/cats-social/internal/repositories"
	"github.com/oklog/ulid/v2"
)

type MatchService interface {
	CreateMatch(ctx context.Context, userId string, payload *matchentity.CreateMatchRequest) error
	GetMatches(ctx context.Context, userId string) ([]*matchentity.GetMatchResponse, error)
	ApproveMatch(ctx context.Context, userId string, payload *matchentity.ApproveMatchRequest) error
	RejectMatch(ctx context.Context, userId string, payload *matchentity.RejectMatchRequest) error
	DeleteMatch(ctx context.Context, userId, matchId string) error
}

type MatchServiceImpl struct {
	MatchRepository repositories.MatchRepository
	CatRepository   repositories.CatRepository
	UserRepository  repositories.UserRepository
}

func NewMatchService(
	matchRepository repositories.MatchRepository,
	catRepository repositories.CatRepository,
	userRepository repositories.UserRepository,
) MatchService {
	return &MatchServiceImpl{
		MatchRepository: matchRepository,
		CatRepository:   catRepository,
		UserRepository:  userRepository,
	}
}

func (s *MatchServiceImpl) CreateMatch(ctx context.Context, userId string, payload *matchentity.CreateMatchRequest) error {
	isMatchCatIdExists, err := s.CatRepository.VerifyId(ctx, payload.MatchCatId)
	if err != nil {
		return err
	}
	if !isMatchCatIdExists {
		return matcherror.ErrMatchCatIdNotFound
	}

	isUserCatIdExists, err := s.CatRepository.VerifyId(ctx, payload.UserCatId)
	if err != nil {
		return err
	}
	if !isUserCatIdExists {
		return matcherror.ErrUserCatIdNotFound
	}

	isUserCatIdBelongToTheUser, err := s.CatRepository.VerifyOwner(ctx, payload.UserCatId, userId)
	if err != nil {
		return err
	}
	if !isUserCatIdBelongToTheUser {
		return matcherror.ErrUserCatIdNotBelongToTheUser
	}

	isBothCatsHaveSameGender, err := s.MatchRepository.VerifyGenderOfBothCats(ctx, payload.MatchCatId, payload.UserCatId)
	if err != nil {
		return err
	}
	if isBothCatsHaveSameGender {
		return matcherror.ErrBothCatsHaveSameGender
	}

	err = s.MatchRepository.VerifyCatsNotMatched(ctx, payload.MatchCatId, payload.UserCatId)
	if err != nil {
		return err
	}

	isBothCatsHaveSameOwner, err := s.MatchRepository.VerifyOwnerOfBothCats(ctx, payload.MatchCatId, payload.UserCatId)
	if err != nil {
		return err
	}
	if isBothCatsHaveSameOwner {
		return matcherror.ErrBothCatsHaveSameOwner
	}

	isMatchRequestExists, err := s.MatchRepository.VerifyMatchRequestExistence(ctx, payload.MatchCatId, payload.UserCatId)
	if err != nil {
		return err
	}
	if isMatchRequestExists {
		return matcherror.ErrDuplicateMatchRequest
	}

	matchCat := &matchentity.Match{
		Id:         ulid.Make().String(),
		MatchCatId: payload.MatchCatId,
		UserCatId:  payload.UserCatId,
		Message:    payload.Message,
	}

	err = s.MatchRepository.CreateMatch(ctx, matchCat)
	if err != nil {
		return err
	}

	return nil
}

func (s *MatchServiceImpl) GetMatches(ctx context.Context, userId string) ([]*matchentity.GetMatchResponse, error) {
	return s.MatchRepository.GetMatches(ctx, userId)
}

func (s *MatchServiceImpl) ApproveMatch(ctx context.Context, userId string, payload *matchentity.ApproveMatchRequest) error {
	isMatchIdExists, err := s.MatchRepository.VerifyMatchId(ctx, payload.MatchId)
	if err != nil {
		return err
	}
	if !isMatchIdExists {
		return matcherror.ErrMatchIdNotFound
	}

	isMatchIdValid, err := s.MatchRepository.VerifyMatchIdValidity(ctx, payload.MatchId)
	if err != nil {
		return err
	}
	if !isMatchIdValid {
		return matcherror.ErrMatchIdIsNoLongerValid
	}

	isMatchIssuer, err := s.MatchRepository.VerifyMatchIssuer(ctx, payload.MatchId, userId)
	if err != nil {
		return err
	}
	if isMatchIssuer {
		return matcherror.ErrIssuerCannotDecide
	}

	err = s.MatchRepository.ApproveMatch(ctx, payload.MatchId)
	if err != nil {
		return err
	}

	return nil
}

func (s *MatchServiceImpl) RejectMatch(ctx context.Context, userId string, payload *matchentity.RejectMatchRequest) error {
	isMatchIdExists, err := s.MatchRepository.VerifyMatchId(ctx, payload.MatchId)
	if err != nil {
		return err
	}
	if !isMatchIdExists {
		return matcherror.ErrMatchIdNotFound
	}

	isMatchIdValid, err := s.MatchRepository.VerifyMatchIdValidity(ctx, payload.MatchId)
	if err != nil {
		return err
	}
	if !isMatchIdValid {
		return matcherror.ErrMatchIdIsNoLongerValid
	}

	isMatchIssuer, err := s.MatchRepository.VerifyMatchIssuer(ctx, payload.MatchId, userId)
	if err != nil {
		return err
	}
	if isMatchIssuer {
		return matcherror.ErrIssuerCannotDecide
	}

	err = s.MatchRepository.RejectMatch(ctx, payload.MatchId)
	if err != nil {
		return err
	}

	return nil
}

func (s *MatchServiceImpl) DeleteMatch(ctx context.Context, userId, matchId string) error {
	isMatchIssuer, err := s.MatchRepository.VerifyMatchIssuer(ctx, matchId, userId)
	if err != nil {
		return err
	}
	if !isMatchIssuer {
		return matcherror.ErrNotIssuer
	}

	isMatchIdExists, err := s.MatchRepository.VerifyMatchId(ctx, matchId)
	if err != nil {
		return err
	}
	if !isMatchIdExists {
		return matcherror.ErrMatchIdNotFound
	}

	isMatchIdValid, err := s.MatchRepository.VerifyMatchIdValidity(ctx, matchId)
	if err != nil {
		return err
	}
	if !isMatchIdValid {
		return matcherror.ErrMatchIdIsNoLongerValid
	}

	err = s.MatchRepository.DeleteMatch(ctx, matchId)
	if err != nil {
		return err
	}

	return nil
}
