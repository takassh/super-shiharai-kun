package util

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID    string `json:"user_id"`
	CompanyID string `json:"company_id"`
	jwt.RegisteredClaims
}
