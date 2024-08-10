package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
	"github.com/morozoffnor/gophermart-diploma/internal/config"
)

type ContextUserIDKey string

var ContextUserID ContextUserIDKey = "user_id"

type Auth struct {
	cfg *config.Config
	Jwt *JWT
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

type User struct {
	Id            string `json:"userID"`
	Login         string `json:"login"`
	Password      string `json:"password"`
	Authenticated bool
}

func New(cfg *config.Config) *Auth {
	return &Auth{cfg: cfg, Jwt: &JWT{secret: "supersecret"}}
}

func (a *Auth) HashPassword(p string) string {
	h := sha256.New()
	h.Write([]byte(p))
	return hex.EncodeToString(h.Sum(nil))
}
