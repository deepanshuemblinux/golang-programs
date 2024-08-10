package user

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}
