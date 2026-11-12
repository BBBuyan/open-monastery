package token

import (
	"errors"
	"strconv"
	"time"

	"example.com/open-monastery/internal/app/abstractions"
	"github.com/golang-jwt/jwt/v5"
)

type JwtIssuer struct {
	signingKey    []byte
	Audiences     []string
	SigningMethod jwt.SigningMethod
	Issuer        string
	ValidMinute   int
}

var _ abstractions.TokeIssuer = (*JwtIssuer)(nil)

func NewJwtIssuer(signingKey []byte) *JwtIssuer {
	return &JwtIssuer{
		signingKey:    signingKey,
		SigningMethod: jwt.SigningMethodHS256,
		Issuer:        "open-monastery",
		Audiences:     []string{"open-monastery-api"},
	}
}

func (j *JwtIssuer) IssueAccessToken(userId int64) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:  strconv.FormatInt(userId, 10),
		Issuer:   j.Issuer,
		Audience: jwt.ClaimStrings(j.Audiences),
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(
			now.Add(time.Duration(j.ValidMinute) * time.Minute),
		),
	}

	return jwt.NewWithClaims(j.SigningMethod, claims).SignedString(j.signingKey)
}

func (j *JwtIssuer) ValidateToken(tokenStr string) (*string, error) {
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&claims,
		func(token *jwt.Token) (any, error) {
			return j.signingKey, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(j.Issuer),
		jwt.WithAudience(j.Audiences...),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid || claims.Subject == "" {
		return nil, errors.New("invalid token")
	}

	return &claims.Subject, nil
}
