package main

import (
	_ "example.com/open-monastery/docs"

	"example.com/open-monastery/internal/book"
	"example.com/open-monastery/internal/handlers"
	"example.com/open-monastery/internal/monastery"
	"example.com/open-monastery/internal/mysql"
	"example.com/open-monastery/internal/reading"
	"example.com/open-monastery/internal/user"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Open Monastery API
// @version 1.0
// @description API for Open Monastery
// @BasePath /api/v1
func main() {
	db, err := mysql.Open("root", "Password123!", "localhost", "3306", "om")
	if err != nil {
		panic(err)
	}
	bookRepository := mysql.NewBookRepo(db)
	bookService := book.NewService(bookRepository)
	bookHandler := handlers.NewBookHandler(bookService)

	monasteryRepository := mysql.NewMonasteryRepo(db)
	monasteryService := monastery.NewService(monasteryRepository)
	monasteryHandler := handlers.NewMonasteryHandler(monasteryService)

	userRepository := mysql.NewUserRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	readingRepository := mysql.NewReadingRepository(db)
	readingService := reading.NewService(
		readingRepository,
		userRepository,
		bookRepository,
	)
	readingHandler := handlers.NewReadingHandler(readingService)

	r := gin.Default()
	r.Use(handlers.ErrorHandler())
	api := r.Group("/api/v1")

	books := api.Group("/books")
	{
		books.GET("/:id", bookHandler.GetBook)
		books.POST("", bookHandler.AddBook)
	}

	monasteries := api.Group("/monasteries")
	{
		monasteries.GET("/:id", monasteryHandler.GetMonastery)
		monasteries.POST("", monasteryHandler.AddMonastery)
	}

	users := api.Group("/users")
	{
		users.GET("/:id", userHandler.GetUser)
		users.POST("", userHandler.AddUser)
	}

	readings := api.Group("/readings")
	{
		readings.GET("/:id", readingHandler.GetReading)
		readings.POST("", readingHandler.AddReading)
	}

	r.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	r.Run("127.0.0.1:8080")
}
