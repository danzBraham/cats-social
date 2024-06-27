package matcherror

import (
	"errors"
)

var (
	ErrMatchIdNotFound             = errors.New("match id not found")
	ErrMatchIdIsNoLongerValid      = errors.New("match id is no longer valid")
	ErrIssuerCannotDecide          = errors.New("match issuer can't make the decision")
	ErrNotIssuer                   = errors.New("you are not the issuer")
	ErrMatchCatIdNotFound          = errors.New("match cat id not found")
	ErrUserCatIdNotFound           = errors.New("user cat id not found")
	ErrUserCatIdNotBelongToTheUser = errors.New("user cat id not belong to the user")
	ErrBothCatsHaveSameGender      = errors.New("both cats have same gender")
	ErrBothCatsHaveAlreadyMatched  = errors.New("both cats have already matched")
	ErrBothCatsHaveSameOwner       = errors.New("both cats have same owner")
	ErrDuplicateMatchRequest       = errors.New("match request already exists for these two cats")
)
