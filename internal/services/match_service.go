package services

import (
	"context"

	"github.com/danzBraham/cats-social/internal/entities/catentity"
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

	isBothCatsHaveSameGender, err := s.MatchRepository.VerifyBothCatsGender(ctx, payload.MatchCatId, payload.UserCatId)
	if err != nil {
		return err
	}
	if isBothCatsHaveSameGender {
		return matcherror.ErrBothCatsHaveSameGender
	}

	err = s.MatchRepository.VerifyBothCatsNotMatched(ctx, payload.MatchCatId, payload.UserCatId)
	if err != nil {
		return err
	}

	isBothCatsHaveTheSameOwner, err := s.MatchRepository.VerifyBothCatsHaveTheSameOwner(ctx, payload.MatchCatId, payload.UserCatId)
	if err != nil {
		return err
	}
	if isBothCatsHaveTheSameOwner {
		return matcherror.ErrBothCatsHaveSameOwner
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
	matches, err := s.MatchRepository.GetMatches(ctx, userId)
	if err != nil {
		return nil, err
	}

	matchResponses := []*matchentity.GetMatchResponse{}
	for _, match := range matches {
		matchCatDetail, err := s.CatRepository.GetCatById(ctx, match.MatchCatId)
		if err != nil {
			return nil, err
		}

		userCatDetail, err := s.CatRepository.GetCatById(ctx, match.UserCatId)
		if err != nil {
			return nil, err
		}

		issuerDetail, err := s.UserRepository.GetUserById(ctx, userCatDetail.OwnerId)
		if err != nil {
			return nil, err
		}

		matchResponses = append(matchResponses, &matchentity.GetMatchResponse{
			Id: match.Id,
			IssuedBy: matchentity.IssuerDetail{
				Name:      issuerDetail.Name,
				Email:     issuerDetail.Email,
				CreatedAt: issuerDetail.CreatedAt,
			},
			MatchCatDetail: catentity.GetCatResponse{
				Id:          matchCatDetail.Id,
				Name:        matchCatDetail.Name,
				Race:        matchCatDetail.Race,
				Sex:         matchCatDetail.Sex,
				AgeInMonth:  matchCatDetail.AgeInMonth,
				Description: matchCatDetail.Description,
				ImageUrls:   matchCatDetail.ImageUrls,
				HasMatched:  matchCatDetail.HasMatched,
				CreatedAt:   matchCatDetail.CreatedAt,
			},
			UserCatDetail: catentity.GetCatResponse{
				Id:          userCatDetail.Id,
				Name:        userCatDetail.Name,
				Race:        userCatDetail.Race,
				Sex:         userCatDetail.Sex,
				AgeInMonth:  userCatDetail.AgeInMonth,
				Description: userCatDetail.Description,
				ImageUrls:   userCatDetail.ImageUrls,
				HasMatched:  userCatDetail.HasMatched,
				CreatedAt:   userCatDetail.CreatedAt,
			},
			Message:   match.Message,
			CreatedAt: match.CreatedAt,
		})
	}

	return matchResponses, nil
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
