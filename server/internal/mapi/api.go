package mapi

import (
	"database/sql"

	_ "example.com/open-monastery/docs"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/services"
	"example.com/open-monastery/internal/domain/models"
	database "example.com/open-monastery/internal/infra/db"
	"example.com/open-monastery/internal/infra/password"
	"example.com/open-monastery/internal/mapi/handlers"
	"example.com/open-monastery/internal/mapi/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func addBook(
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
func addCategory(
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

func addMonastery(
	db *sql.DB,
	api *gin.RouterGroup,
) {
	monasteryRepository := database.NewMonasteryRepo(db)
	monasteryService := services.NewMonasteryService(monasteryRepository)
	monasteryHandler := handlers.NewMonasteryHandler(monasteryService)

	monasteries := api.Group("/monasteries")
	{
		monasteries.GET("/:id", withProblemDetails(monasteryHandler.GetMonastery))
		monasteries.GET("", withProblemDetails(monasteryHandler.GetMonasteries))
		monasteries.POST("",
			withProblemDetails(middlewares.RequireRole(models.Admin)),
			withProblemDetails(monasteryHandler.AddMonastery),
		)
	}
}

func addReadings(db *sql.DB, authApi *gin.RouterGroup) {
	readingRepository := database.NewReadingRepository(db)
	readingService := services.NewReadingService(readingRepository)
	readingHandler := handlers.NewReadingHandler(readingService)
	readings := authApi.Group("/readings")
	{
		readings.GET("", withProblemDetails(readingHandler.GetReadings))
		readings.POST("", withProblemDetails(readingHandler.AddReading))
		readings.POST("/:id/claim",
			withProblemDetails(middlewares.RequireRole(models.Monk)),
			withProblemDetails(readingHandler.ClaimReading),
		)
	}
}

func addSwagger(
	r *gin.Engine,
) {
	r.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
}
func addAuth(
	db *sql.DB,
	api *gin.RouterGroup,
	tokenIssuer abstractions.TokeIssuer,
) {
	authRepo := database.NewAuthRepository(db)
	passHasher := password.NewArgonPasswordHasher()
	authService := services.NewAuthService(authRepo, passHasher, tokenIssuer)
	authHandler := handlers.NewAuthHandler(authService)
	api.POST("/auth/login", withProblemDetails(authHandler.Login))
	api.POST("/auth/register/monks", withProblemDetails(authHandler.RegisterMonk))
	api.POST("/auth/register/users", withProblemDetails(authHandler.RegisterUser))
}

func RegisterAPI(
	db *sql.DB,
	tokenIssuer abstractions.TokeIssuer,
) *gin.Engine {
	r := gin.Default()
	authFunc := middlewares.RequireAuth(tokenIssuer.ValidateToken)

	api := r.Group("/api/v1")
	authApi := api.Group("", authFunc)

	addAuth(db, api, tokenIssuer)
	addBook(db, api)
	addCategory(db, api)
	addReadings(db, authApi)
	addMonastery(db, api)

	addSwagger(r)

	return r
}
