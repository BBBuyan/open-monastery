package main

import (
	"example.com/open-monastery/internal/infra/db"
	"example.com/open-monastery/internal/mapi"
)

// @title Open Monastery API
// @version 1.0
// @description API for Open Monastery
// @BasePath /api/v1
func main() {
	db, err := db.Open("root", "Password123!", "localhost", "3306", "om")
	if err != nil {
		panic(err)
	}
	r := mapi.RegisterAPI(db)
	r.Run("127.0.0.1:8080")
}
