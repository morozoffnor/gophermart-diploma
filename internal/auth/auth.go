package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/morozoffnor/gophermart-diploma/internal/config"
)

type ContextUserIDKey string

var ContextUserID ContextUserIDKey = "user_id"

type Auth struct {
	cfg *config.Config
	jwt *JWT
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

type User struct {
	Id            string `json:"userID"`
	Login         string `json:"login"`
	Password      string `json:"password"`
	authenticated bool
}

func New(cfg *config.Config) *Auth {
	return &Auth{cfg: cfg, jwt: &JWT{secret: "supersecret"}}
}

func (a *Auth) Register(u *User) (bool, error) {

}
