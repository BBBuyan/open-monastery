package token

import (
	"strconv"
	"time"

	"example.com/open-monastery/internal/app/abstractions"
	"github.com/golang-jwt/jwt/v5"
)

type JwtIssuer struct {
	SigningKey    []byte
	SigningMethod jwt.SigningMethod
}

var _ abstractions.TokeIssuer = (*JwtIssuer)(nil)

func NewJwtIssuer(signingKey []byte) *JwtIssuer {
	return &JwtIssuer{
		SigningKey:    signingKey,
		SigningMethod: jwt.SigningMethodHS256,
	}
}

func (j *JwtIssuer) IssueAccessToken(userId int64) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userId, 10),
		Issuer:    "open-monastery",
		Audience:  jwt.ClaimStrings{"open-monastery-api"},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
	}

	return jwt.NewWithClaims(j.SigningMethod, claims).SignedString(j.SigningKey)
}
