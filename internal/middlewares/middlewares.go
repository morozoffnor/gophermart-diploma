package middlewares

import (
	"context"
	"github.com/morozoffnor/gophermart-diploma/internal/auth"
	"log"
	"net/http"
	"time"
)

func NewAuthMiddleware(a *auth.Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("Authorization")
			if err != nil {
				log.Print(err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := a.Jwt.ParseToken(cookie.Value)
			if err != nil {
				log.Print(err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// обновить токен, если скоро стухнет
			if claims.ExpiresAt.Add(time.Hour).After(time.Now()) {
				token, _ := a.Jwt.GenerateToken(claims.UserID)
				ctx, _ := a.Jwt.AddTokenToCookies(&w, r, token)
				r = r.WithContext(ctx)
			}

			ctx := context.WithValue(r.Context(), auth.ContextUserID, claims.UserID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
