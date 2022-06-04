package middlewares

import (
	"context"
	"go-web-api/internal/protocols"
	"net/http"
)

type key int

var UserKey key

type jwtProvider interface {
	Validate(token string) (*int, error)
}

func AuthMiddleware(p jwtProvider, publicPaths map[string]any) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if _, ok := publicPaths[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
			} else {

				token := r.Header.Get("Authorization")

				if len(token) == 0 {
					protocols.Unauthorized(w)
					return
				}

				userId, err := p.Validate(token)

				if err != nil {
					protocols.Unauthorized(w)
					return
				}

				ctx := r.Context()
				userCtx := context.WithValue(ctx, UserKey, userId)
				req := r.WithContext(userCtx)

				next.ServeHTTP(w, req)
			}
		})
	}
}
