package securities_impl

import (
	"fmt"
	"time"

	"github.com/danzbraham/cats-social/internal/applications/securities"
	auth_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/auth"
	auth_entity "github.com/danzbraham/cats-social/internal/domains/entities/auth"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenManager struct {
	Key []byte
}

func NewJWTTokenManager(key []byte) securities.AuthTokenManager {
	return &JWTTokenManager{Key: key}
}

type CustomClaims struct {
	UserId string
	jwt.RegisteredClaims
}

func (j *JWTTokenManager) GenerateToken(ttl time.Duration, userId string) (string, error) {
	now := time.Now()
	expiry := now.Add(ttl)

	claims := &CustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Key)
}

func (j *JWTTokenManager) VerifyToken(tokenString string) (*auth_entity.Credential, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
		return j.Key, nil
	})
	if token == nil {
		return nil, auth_exception.ErrMissingToken
	}
	if err != nil {
		return nil, auth_exception.ErrInvalidToken
	}
	if !token.Valid {
		return nil, auth_exception.ErrInvalidToken
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || claims == nil {
		return nil, auth_exception.ErrUnknownClaims
	}

	return &auth_entity.Credential{
		UserId: claims.UserId,
	}, nil
}
