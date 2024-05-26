package match_exception

import "errors"

var (
	ErrMatchIdIsNotFound         = errors.New("match id is not found")
	ErrMatchIdIsNoLongerValid    = errors.New("match id is no longer valid")
	ErrMatchCatIsNotFound        = errors.New("match cat is not found")
	ErrMatchCatIdIsNotFound      = errors.New("match cat id is not found")
	ErrUserCatIdIsNotFound       = errors.New("user cat id is not found")
	ErrUserCatIdNotBelongTheUser = errors.New("user cat id is not belong to the user")
	ErrBothCatHaveSameGender     = errors.New("both cats have the same gender")
	ErrBothCatAlreadyMatched     = errors.New("both cats already matched")
	ErrBothCatsHaveTheSameOwner  = errors.New("both cats have the same owner")
)
