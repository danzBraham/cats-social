package securities

import (
	"time"

	auth_entity "github.com/danzbraham/cats-social/internal/domains/entities/auth"
)

type AuthTokenManager interface {
	GenerateToken(ttl time.Duration, userID string) (token string, err error)
	VerifyToken(tokenString string) (*auth_entity.Credential, error)
}
