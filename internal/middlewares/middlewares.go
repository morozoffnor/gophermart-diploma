package middlewares

import (
	"context"
	"github.com/morozoffnor/gophermart-diploma/internal/auth"
	"github.com/morozoffnor/gophermart-diploma/internal/storage"
	"log"
	"net/http"
	"time"
)

type Middlewares struct {
	auth *auth.Auth
	db   *storage.DB
}

func New(a *auth.Auth, db *storage.DB) *Middlewares {
	return &Middlewares{
		auth: a,
		db:   db,
	}
}

func (m *Middlewares) Auth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("Authorization")
			if err != nil {
				log.Print(err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}

			claims, err := m.auth.Jwt.ParseToken(cookie.Value)
			if err != nil {
				log.Print(err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}

			// обновить токен, если скоро стухнет
			if claims.ExpiresAt.Add(time.Hour).After(time.Now()) {
				token, _ := m.auth.Jwt.GenerateToken(claims.UserID)
				ctx, _ := m.auth.Jwt.AddTokenToCookies(&w, r, token)
				r = r.WithContext(ctx)
			}

			ctx := context.WithValue(r.Context(), auth.ContextUserID, claims.UserID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
