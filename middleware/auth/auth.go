package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type ctxKey string

var (
	authCtxKey ctxKey = "auth-ctx"
)

type AuthContext struct {
	Subject string
}

func WithAuthContext(parent context.Context, authCtx *AuthContext) context.Context {
	return context.WithValue(parent, authCtxKey, authCtx)
}

func AuthContextFromContext(ctx context.Context) (*AuthContext, bool) {
	authCtx, ok := ctx.Value(authCtxKey).(*AuthContext)
	if !ok {
		return nil, false
	}

	return authCtx, true
}

func JWTTokenVerifierMiddleware(secret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

			if len(authHeader) != 2 {
				http.Error(w, "malformed token", http.StatusUnauthorized)
				return
			}

			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unsupported signing method: %v", token.Header["alg"])
				}

				return []byte(secret), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			var ctx context.Context
			if ok && token.Valid {
				sub := claims["sub"].(string)
				ctx = WithAuthContext(r.Context(), &AuthContext{Subject: sub})
			} else {
				http.Error(w, "malformed token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}