package main

import (
	"example.com/open-monastery/internal/api"
	"example.com/open-monastery/internal/infra/db"
	"example.com/open-monastery/internal/infra/token"
)

// devSigningKey is for local development only.
const devSigningKey = "open-monastery-local-development-signing-key"

// main
//
//	@title						Open Monastery API
//	@version					1.0
//	@description				API for Open Monastery
//	@BasePath					/api/v1
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
func main() {
	jwtIssuer := token.NewJwtIssuer([]byte(devSigningKey))
	jwtIssuer.ValidMinute = 15

	db, err := db.Open("root", "Password123!", "localhost", "3306", "om")
	if err != nil {
		panic(err)
	}
	r := api.RegisterAPI(db, jwtIssuer)
	r.Run("127.0.0.1:8080")
}
