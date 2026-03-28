package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"ivorioci-stream-service/models"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type UserContext struct {
	Sub   string
	Email string
}

func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, appErr := extractUser(r, secret)
			if appErr != nil {
				writeError(w, appErr.StatusCode, appErr.Code, appErr.Message)
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(r *http.Request) (*UserContext, bool) {
	user, ok := r.Context().Value(UserContextKey).(*UserContext)
	return user, ok
}

func extractUser(r *http.Request, secret string) (*UserContext, *models.AppError) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, models.ErrUnauthorized
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrTokenInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return nil, models.ErrTokenExpired
		}
		return nil, models.ErrTokenInvalid
	}

	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		return nil, models.ErrTokenInvalid
	}

	return &UserContext{Sub: sub, Email: claims.Email}, nil
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.NewError(code, message))
}
