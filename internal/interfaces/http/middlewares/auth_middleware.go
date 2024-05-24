package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/danzbraham/cats-social/internal/applications/securities"
	http_common "github.com/danzbraham/cats-social/internal/commons/http"
)

type ContextKey string

var ContextUserIdKey ContextKey = "userId"

type AuthMiddleware struct {
	AuthTokenManager securities.AuthTokenManager
}

func NewAuthMiddleware(authTokenManager securities.AuthTokenManager) *AuthMiddleware {
	return &AuthMiddleware{AuthTokenManager: authTokenManager}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http_common.ResponseError(w, http.StatusUnauthorized, "Unauthorized error", "Missing Authorizaton header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http_common.ResponseError(w, http.StatusUnauthorized, "Unauthorized error", "Invalid Authorization header format")
			return
		}

		credential, err := m.AuthTokenManager.VerifyToken(tokenString)
		if err != nil {
			http_common.ResponseError(w, http.StatusUnauthorized, "Unauthorized error", err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserIdKey, credential.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
