package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"main/config"
	"time"
)

type TokenClaim struct {
	UserId       int64  `json:"user_id,omitempty"`
	Lang         string `json:"lang,omitempty"`
	UserType     int8   `json:"user_type,omitempty"`
	PersonalCode string `json:"personal_code,omitempty"`
	UserTitle    string `json:"user_title,omitempty"`
	UserStatus   int8   `json:"user_status,omitempty"`
	Exp          int64  `json:"exp"`
	Iat          int64  `json:"iat"`
	jwt.RegisteredClaims
}

func GenerateToken(cfg *config.Config, claims TokenClaim) (string, error) {
	signingKey := []byte(cfg.Server.APP_SECRET)
	claims.Iat = time.Now().Unix()
	claims.Exp = time.Now().Local().Add(time.Hour * time.Duration(cfg.Server.JWT_TOKEN_EXPIRE_TIME)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}
