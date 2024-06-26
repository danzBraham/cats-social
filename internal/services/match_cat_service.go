package services

import (
	"context"

	"github.com/danzBraham/cats-social/internal/entities/catentity"
	"github.com/danzBraham/cats-social/internal/entities/matchcatentity"
	"github.com/danzBraham/cats-social/internal/errors/matchcaterror"
	"github.com/danzBraham/cats-social/internal/repositories"
	"github.com/oklog/ulid/v2"
)

type MatchCatService interface {
	CreateMatchCat(ctx context.Context, userId string, payload *matchcatentity.CreateMatchCatRequest) error
	GetMatchCats(ctx context.Context, userId string) ([]*matchcatentity.GetMatchCatResponse, error)
	ApproveMatchCat(ctx context.Context, payload *matchcatentity.ApproveMatchCatRequest) error
}

type MatchCatServiceImpl struct {
	MatchCatRepository repositories.MatchCatRepository
	CatRepository      repositories.CatRepository
	UserRepository     repositories.UserRepository
}

func NewMatchCatService(
	matchCatRepository repositories.MatchCatRepository,
	catRepository repositories.CatRepository,
	userRepository repositories.UserRepository,
) MatchCatService {
	return &MatchCatServiceImpl{
		MatchCatRepository: matchCatRepository,
		CatRepository:      catRepository,
		UserRepository:     userRepository,
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

func (s *MatchCatServiceImpl) GetMatchCats(ctx context.Context, userId string) ([]*matchcatentity.GetMatchCatResponse, error) {
	matchCats, err := s.MatchCatRepository.GetMatchCats(ctx, userId)
	if err != nil {
		return nil, err
	}

	matchCatResponses := []*matchcatentity.GetMatchCatResponse{}
	for _, matchCat := range matchCats {
		issuerDetail, err := s.UserRepository.GetUserById(ctx, matchCat.IssuedBy)
		if err != nil {
			return nil, err
		}

		matchCatDetail, err := s.CatRepository.GetCatById(ctx, matchCat.MatchCatId)
		if err != nil {
			return nil, err
		}

		userCatDetail, err := s.CatRepository.GetCatById(ctx, matchCat.UserCatId)
		if err != nil {
			return nil, err
		}

		matchCatResponses = append(matchCatResponses, &matchcatentity.GetMatchCatResponse{
			Id: matchCat.Id,
			IssuedBy: matchcatentity.IssuerDetail{
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
			Message:   matchCat.Message,
			CreatedAt: matchCat.CreatedAt,
		})
	}

	return matchCatResponses, nil
}

func (s *MatchCatServiceImpl) ApproveMatchCat(ctx context.Context, payload *matchcatentity.ApproveMatchCatRequest) error {
	isMatchIdExists, err := s.MatchCatRepository.VerifyMatchId(ctx, payload.MatchId)
	if err != nil {
		return err
	}
	if !isMatchIdExists {
		return matchcaterror.ErrMatchIdNotFound
	}

	isMatchIdValid, err := s.MatchCatRepository.VerifyMatchIdValidity(ctx, payload.MatchId)
	if err != nil {
		return err
	}
	if !isMatchIdValid {
		return matchcaterror.ErrMatchIdIsNoLongerValid
	}

	err = s.MatchCatRepository.ApproveMatchCat(ctx, payload.MatchId)
	if err != nil {
		return err
	}

	return nil
}
