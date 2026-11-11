package mapi

import (
	"database/sql"

	_ "example.com/open-monastery/docs"

	"example.com/open-monastery/internal/app/services"
	database "example.com/open-monastery/internal/infra/db"
	"example.com/open-monastery/internal/infra/password"
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
		books.GET("/:id", withProblemDetails(bookHandler.GetBook))
		books.GET("", withProblemDetails(bookHandler.GetBooks))
		books.POST("", withProblemDetails(bookHandler.AddBook))
		books.GET("/:id/categories", withProblemDetails(bookHandler.GetBookCategories))
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
		categories.GET("", withProblemDetails(categoryHandler.GetCategories))
		categories.POST("", withProblemDetails(categoryHandler.AddCategory))
		categories.POST("/:id/books", withProblemDetails(categoryHandler.AddBook))
		categories.GET("/:id/books", withProblemDetails(categoryHandler.GetBooksByCategoryId))
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
		monasteries.GET("/:id", withProblemDetails(monasteryHandler.GetMonastery))
		monasteries.POST("", withProblemDetails(monasteryHandler.AddMonastery))
	}
}

func registerUser(
	db *sql.DB,
	api *gin.RouterGroup,
) {
	userRepository := database.NewUserRepository(db)
	passwordHasher := password.NewArgonPasswordHasher()
	userService := services.NewUserService(userRepository, passwordHasher)
	userHandler := handlers.NewUserHandler(userService)
	api.POST("/monks", withProblemDetails(userHandler.RegisterMonk))

	users := api.Group("/users")
	{
		users.GET("/:id", withProblemDetails(userHandler.GetUser))
		users.POST("", withProblemDetails(userHandler.RegisterUser))
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
		readings.GET("/:id", withProblemDetails(readingHandler.GetReading))
		readings.POST("", withProblemDetails(readingHandler.AddReading))
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
	api := r.Group("/api/v1")

	registerBook(db, api)
	registerCategory(db, api)
	registerUser(db, api)
	registerReadings(db, api)
	registerMonastery(db, api)

	registerSwagger(r)

	return r
}
