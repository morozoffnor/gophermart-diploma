package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type JWT struct {
	secret string
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func (j *JWT) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UserID: userID,
	})
	tokenStr, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return "Bearer " + tokenStr, nil
}

func (j *JWT) ParseToken(token string) (*Claims, error) {
	token = token[7:]
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (j *JWT) CheckToken(r *http.Request) bool {
	cookie, err := r.Cookie("Authorization")

	if err != nil {
		return false
	}

	claims, err := j.ParseToken(cookie.Value)
	if err != nil {
		return false
	}
	if claims.UserID == "" {
		return false
	}
	return true
}

func (j *JWT) GetUserIdFromToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("Authorization")

	if err != nil {
		return "", err
	}

	claims, err := j.ParseToken(cookie.Value)
	if err != nil {
		return "", err
	}
	if claims.UserID == "" {
		return "", err
	}
	return claims.UserID, nil
}

func (j *JWT) AddTokenToCookies(w *http.ResponseWriter, r *http.Request, token string) (context.Context, error) {
	http.SetCookie(*w, &http.Cookie{
		Name:    "Authorization",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})

	claims, err := j.ParseToken(token)
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)

	return ctx, nil
}
