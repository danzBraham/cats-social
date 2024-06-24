package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/danzBraham/cats-social/internal/errors/autherror"
	"github.com/danzBraham/cats-social/internal/helpers/http_helper"
	"github.com/danzBraham/cats-social/internal/helpers/jwt"
)

type ContextKey string

var ContextUserIdKey ContextKey = "userId"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http_helper.HandleErrorResponse(w, http.StatusUnauthorized, autherror.ErrMissingAuthHeader)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http_helper.HandleErrorResponse(w, http.StatusUnauthorized, autherror.ErrInvalidAuthHeader)
			return
		}

		token, err := jwt.VerifyToken(tokenString)
		if err != nil {
			http_helper.HandleErrorResponse(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserIdKey, token.UserId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
