package caterror

import "errors"

var (
	ErrCatIdNotFound = errors.New("cat id not found")
	ErrCatNotFound   = errors.New("cat not found")
	ErrNotCatOwner   = errors.New("you're not the cat owner")
	ErrSexIsEdited   = errors.New("sex cannot be changed when cat is already requested for a match")
)
