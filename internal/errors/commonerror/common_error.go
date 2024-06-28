package commonerror

import "errors"

var (
	ErrRouteDoesNotExist = errors.New("route does not exist")
	ErrMethodNotAllowed  = errors.New("method is not allowed")
)
