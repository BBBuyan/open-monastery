package main

import (
	"example.com/open-monastery/internal/infra/db"
	"example.com/open-monastery/internal/infra/token"
	"example.com/open-monastery/internal/mapi"
)

// devSigningKey is for local development only.
const devSigningKey = "open-monastery-local-development-signing-key"

// @title Open Monastery API
// @version 1.0
// @description API for Open Monastery
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter Bearer followed by a space and the JWT access token.
func main() {
	jwtIssuer := token.NewJwtIssuer([]byte(devSigningKey))
	jwtIssuer.ValidMinute = 15

	db, err := db.Open("root", "Password123!", "localhost", "3306", "om")
	if err != nil {
		panic(err)
	}
	r := mapi.RegisterAPI(db, jwtIssuer)
	r.Run("127.0.0.1:8080")
}
