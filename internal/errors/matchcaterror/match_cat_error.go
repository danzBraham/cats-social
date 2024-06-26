package matchcaterror

import (
	"errors"
	"net/http"
)

var (
	ErrMatchIdNotFound             = errors.New("match id not found")
	ErrMatchIdIsNoLongerValid      = errors.New("match id is no longer valid")
	ErrUnauthorizedDecision        = errors.New("request issuer is not authorized to make a decision")
	ErrMatchCatIdNotFound          = errors.New("match cat id not found")
	ErrUserCatIdNotFound           = errors.New("user cat id not found")
	ErrUserCatIdNotBelongToTheUser = errors.New("user cat id not belong to the user")
	ErrBothCatsHaveTheSameGender   = errors.New("both cats have the same gender")
	ErrBothCatsHaveAlreadyMatched  = errors.New("both cats have already matched")
	ErrBothCatsHaveTheSameOwner    = errors.New("both cats have the same owner")
)

var MatchCatErrorMap = map[error]int{
	ErrMatchIdNotFound:             http.StatusNotFound,
	ErrMatchCatIdNotFound:          http.StatusNotFound,
	ErrUserCatIdNotFound:           http.StatusNotFound,
	ErrUserCatIdNotBelongToTheUser: http.StatusNotFound,
	ErrBothCatsHaveTheSameGender:   http.StatusBadRequest,
	ErrBothCatsHaveAlreadyMatched:  http.StatusBadRequest,
	ErrBothCatsHaveTheSameOwner:    http.StatusBadRequest,
}
