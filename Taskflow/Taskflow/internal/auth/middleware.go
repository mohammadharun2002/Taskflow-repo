package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func (s *Service) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		authorization := r.Header.Get("Authorization")

		parts := strings.Fields(authorization)
		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			writeAuthError(
				w,
				http.StatusUnauthorized,
				ErrInvalidToken.Error(),
			)
			return
		}

		userID, err := s.validateToken(parts[1])
		if err != nil {
			writeAuthError(
				w,
				http.StatusUnauthorized,
				ErrInvalidToken.Error(),
			)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			userIDContextKey,
			userID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Service) validateToken(tokenString string) (int64, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return s.jwtSecrect, nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
		jwt.WithIssuer("taskflow"),
	)
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		return 0, ErrInvalidToken
	}

	return userID, nil
}

func UserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(userIDContextKey).(int64)
	if !ok || userID <= 0 {
		return 0, errors.New("user ID missing from context")
	}

	return userID, nil
}
