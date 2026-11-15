package token

import (
	"errors"
	"strconv"
	"time"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

type JwtIssuer struct {
	signingKey    []byte
	Audiences     []string
	SigningMethod jwt.SigningMethod
	Issuer        string
	ValidMinute   int
}

type accessClaims struct {
	jwt.RegisteredClaims
	Roles []models.UserType `json:"roles"`
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

func (j *JwtIssuer) IssueAccessToken(userID int64, userType models.UserType) (string, error) {
	now := time.Now()
	claims := accessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			Issuer:    j.Issuer,
			Audience:  jwt.ClaimStrings(j.Audiences),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.ValidMinute) * time.Minute)),
		},
		Roles: []models.UserType{userType},
	}

	return jwt.NewWithClaims(j.SigningMethod, claims).SignedString(j.signingKey)
}

func (j *JwtIssuer) ValidateToken(tokenStr string) (abstractions.AccessIdentity, error) {
	var claims accessClaims

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
		return abstractions.AccessIdentity{}, err
	}
	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if !token.Valid || err != nil || userID <= 0 {
		return abstractions.AccessIdentity{}, errors.New("invalid token")
	}

	return abstractions.AccessIdentity{UserID: userID, Roles: claims.Roles}, nil
}
