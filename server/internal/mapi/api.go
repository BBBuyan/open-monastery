package mapi

import (
	"database/sql"
	_ "example.com/open-monastery/docs"

	"example.com/open-monastery/internal/app/services"
	database "example.com/open-monastery/internal/infra/db"
	"example.com/open-monastery/internal/mapi/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerBook(
	db *sql.DB,
	api *gin.RouterGroup,
) {
	bookRepository := database.NewBookRepo(db)
	bookService := services.NewBookService(bookRepository)
	bookHandler := handlers.NewBookHandler(bookService)

	books := api.Group("/books")
	{
		books.GET("/:id", bookHandler.GetBook)
		books.GET("", bookHandler.GetBooks)
		books.POST("", bookHandler.AddBook)
		books.GET("/:id/categories", bookHandler.GetBookCategories)
	}
}
func registerCategory(
	db *sql.DB,
	api *gin.RouterGroup,
) {
	categoryRepository := database.NewCategoryRepo(db)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	categories := api.Group("/categories")
	{
		categories.GET("", categoryHandler.GetCategories)
		categories.POST("", categoryHandler.AddCategory)
		categories.POST("/:id/books", categoryHandler.AddBook)
		categories.GET("/:id/books", categoryHandler.GetBooksByCategoryId)
	}
}

func registerMonastery(
	db *sql.DB,
	api *gin.RouterGroup,
) {
	monasteryRepository := database.NewMonasteryRepo(db)
	monasteryService := services.NewMonasteryService(monasteryRepository)
	monasteryHandler := handlers.NewMonasteryHandler(monasteryService)

	monasteries := api.Group("/monasteries")
	{
		monasteries.GET("/:id", monasteryHandler.GetMonastery)
		monasteries.POST("", monasteryHandler.AddMonastery)
	}
}

func registerUser(
	db *sql.DB,
	api *gin.RouterGroup,
) {
	userRepository := database.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	users := api.Group("/users")
	{
		users.GET("/:id", userHandler.GetUser)
		users.POST("", userHandler.RegisterUser)
	}
}

func registerReadings(
	db *sql.DB,
	api *gin.RouterGroup,
) {
	readingRepository := database.NewReadingRepository(db)
	readingService := services.NewReadingService(readingRepository)
	readingHandler := handlers.NewReadingHandler(readingService)
	readings := api.Group("/readings")
	{
		readings.GET("/:id", readingHandler.GetReading)
		readings.POST("", readingHandler.AddReading)
	}
}

func registerSwagger(
	r *gin.Engine,
) {
	r.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
}

func RegisterAPI(
	db *sql.DB,
) *gin.Engine {
	r := gin.Default()
	r.Use(handlers.ErrorHandler())
	api := r.Group("/api/v1")

	registerBook(db, api)
	registerCategory(db, api)
	registerUser(db, api)
	registerReadings(db, api)
	registerMonastery(db, api)

	registerSwagger(r)

	return r
}
