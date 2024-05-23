package securities_impl

import (
	"time"

	"github.com/danzbraham/cats-social/internal/applications/securities"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenManager struct {
	Key []byte
}

func NewJWTTokenManager(key []byte) securities.AuthTokenManager {
	return &JWTTokenManager{Key: key}
}

type CustomClaims struct {
	UserID string
	jwt.RegisteredClaims
}

func (j *JWTTokenManager) GenerateToken(ttl time.Duration, userID string) (string, error) {
	now := time.Now()
	expiry := now.Add(ttl)

	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Key)
}
