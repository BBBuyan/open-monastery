package main

import (
	_ "example.com/open-monastery/docs"

	"example.com/open-monastery/internal/book"
	"example.com/open-monastery/internal/category"
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
	categoryRepository := mysql.NewCategoryRepo(db)
	categoryService := category.NewService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

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
		books.GET("", bookHandler.GetBooks)
		books.POST("", bookHandler.AddBook)
		books.GET("/:id/categories", categoryHandler.GetCategoriesByBookId)
	}

	categories := api.Group("/categories")
	{
		categories.GET("", categoryHandler.GetCategories)
		categories.POST("", categoryHandler.AddCategory)
		categories.POST("/:id/books", categoryHandler.AddBook)
		categories.GET("/:id/books", categoryHandler.GetBooksByCategoryId)
	}

	monasteries := api.Group("/monasteries")
	{
		monasteries.GET("/:id", monasteryHandler.GetMonastery)
		monasteries.POST("", monasteryHandler.AddMonastery)
	}

	users := api.Group("/users")
	{
		users.GET("/:id", userHandler.GetUser)
		users.POST("", userHandler.RegisterUser)
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
