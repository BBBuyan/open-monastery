package api

import (
	"database/sql"

	_ "example.com/open-monastery/docs"

	"example.com/open-monastery/internal/api/controllers"
	"example.com/open-monastery/internal/api/middlewares"
	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/services"
	"example.com/open-monastery/internal/domain/enums"
	database "example.com/open-monastery/internal/infra/db"
	"example.com/open-monastery/internal/infra/password"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func addBook(
	db *sql.DB,
	api, authAPI *gin.RouterGroup,
) {
	bookRepository := database.NewBookRepo(db)
	bookService := services.NewBookService(bookRepository)
	bookHandler := controllers.NewBookController(bookService)

	books := api.Group("/monasteries/:id/books")
	{
		books.GET("/:bookId", withProblemDetails(bookHandler.GetBook))
		books.GET("", withProblemDetails(bookHandler.GetBooks))
		books.GET("/:bookId/categories", withProblemDetails(bookHandler.GetBookCategories))
	}
	authAPI.POST(
		"/monasteries/:id/books",
		withProblemDetails(middlewares.RequireRole(enums.Manager, enums.Admin)),
		withProblemDetails(bookHandler.AddBook),
	)
}
func addCategory(
	db *sql.DB,
	api, authAPI *gin.RouterGroup,
) {
	categoryRepository := database.NewCategoryRepo(db)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := controllers.NewCategoryController(categoryService)

	categories := api.Group("/monasteries/:id/categories")
	{
		categories.GET("", withProblemDetails(categoryHandler.GetCategories))
		categories.GET("/:categoryId/books", withProblemDetails(categoryHandler.GetBooksByCategoryId))
	}
	authCategories := authAPI.Group("/monasteries/:id/categories",
		withProblemDetails(middlewares.RequireRole(enums.Manager, enums.Admin)))
	authCategories.POST("", withProblemDetails(categoryHandler.AddCategory))
	authCategories.POST("/:categoryId/books", withProblemDetails(categoryHandler.AddBook))
}

func addMonastery(
	db *sql.DB,
	api, authAPI *gin.RouterGroup,
) {
	monasteryService := services.NewMonasteryService(database.NewMonasteryRepo(db))
	monasteryHandler := controllers.NewMonasteryController(monasteryService)
	monasteries := api.Group("/monasteries")
	{
		monasteries.GET("/:id", withProblemDetails(monasteryHandler.GetMonastery))
		monasteries.GET("", withProblemDetails(monasteryHandler.GetMonasteries))
	}
	adminMonasteries := authAPI.Group("/monasteries", withProblemDetails(middlewares.RequireRole(enums.Admin)))
	adminMonasteries.POST("", withProblemDetails(monasteryHandler.AddMonastery))
}

func addReadings(db *sql.DB, authApi *gin.RouterGroup) {
	readingRepository := database.NewReadingRepository(db)
	readingService := services.NewReadingService(readingRepository)
	readingHandler := controllers.NewReadingController(readingService)
	readings := authApi.Group("/monasteries/:id/readings")
	{
		readings.GET("",
			withProblemDetails(middlewares.RequireRole(enums.Monk)),
			withProblemDetails(readingHandler.GetReadings),
		)
		readings.POST("", withProblemDetails(readingHandler.AddReading))
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
	tokenIssuer abstractions.TokenIssuer,
) {
	authRepo := database.NewAuthRepository(db)
	passHasher := password.NewArgonPasswordHasher()
	authService := services.NewAuthService(authRepo, passHasher, tokenIssuer)
	authController := controllers.NewAuthController(authService)

	api.POST("/auth/login", withProblemDetails(authController.Login))
	api.POST("/auth/register", withProblemDetails(authController.RegisterUser))
	api.POST("/auth/promotion/monk", withProblemDetails(authController.PromoteToMonk))
}

func RegisterAPI(
	db *sql.DB,
	tokenIssuer abstractions.TokenIssuer,
) *gin.Engine {
	r := gin.Default()
	authFunc := middlewares.RequireAuth(tokenIssuer.ValidateToken)

	api := r.Group("/api/v1")
	authApi := api.Group("", authFunc)

	addAuth(db, api, tokenIssuer)
	addBook(db, api, authApi)
	addCategory(db, api, authApi)
	addReadings(db, authApi)
	addMonastery(db, api, authApi)

	addSwagger(r)

	return r
}
