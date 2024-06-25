package caterror

import "errors"

var (
	ErrIdNotFound  = errors.New("id not found")
	ErrSexIsEdited = errors.New("sex is edited when cat is already requested to match")
)
