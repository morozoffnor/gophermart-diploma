package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/morozoffnor/gophermart-diploma/internal/config"
)

type ContextUserIDKey string

var ContextUserID ContextUserIDKey = "user_id"

type Auth struct {
	cfg *config.Config
	Jwt *JWT
}

type User struct {
	ID            string `json:"userID"`
	Login         string `json:"login"`
	Password      string `json:"password"`
	Authenticated bool
}

func New(cfg *config.Config) *Auth {
	// мне просто показалось, что было бы прикольно встроить auth и jwt прям в
	//   структуру с хендлерами, чтобы можно было проще пользоваться всеми методами
	return &Auth{cfg: cfg, Jwt: &JWT{secret: "supersecret"}}
}

func (a *Auth) HashPassword(p string) string {
	h := sha256.New()
	h.Write([]byte(p))
	return hex.EncodeToString(h.Sum(nil))
}
