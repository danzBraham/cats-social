package matcherror

import (
	"errors"
)

var (
	ErrMatchIdNotFound             = errors.New("match id not found")
	ErrMatchIdIsNoLongerValid      = errors.New("match id is no longer valid")
	ErrUnauthorizedDecision        = errors.New("request issuer is not authorized to make a decision")
	ErrMatchCatIdNotFound          = errors.New("match cat id not found")
	ErrUserCatIdNotFound           = errors.New("user cat id not found")
	ErrUserCatIdNotBelongToTheUser = errors.New("user cat id not belong to the user")
	ErrBothCatsHaveSameGender      = errors.New("both cats have same gender")
	ErrBothCatsHaveAlreadyMatched  = errors.New("both cats have already matched")
	ErrBothCatsHaveSameOwner       = errors.New("both cats have same owner")
)
