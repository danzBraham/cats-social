package securities

import "time"

type AuthTokenManager interface {
	GenerateToken(ttl time.Duration, userID string) (token string, err error)
}
