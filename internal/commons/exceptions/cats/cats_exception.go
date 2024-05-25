package cat_exception

import "errors"

var (
	ErrCatNotFound        = errors.New("cat not found")
	ErrCatIdAlreadyExists = errors.New("cat id already exists")
	ErrCatIdIsNotFound    = errors.New("cat id is not found")
)
