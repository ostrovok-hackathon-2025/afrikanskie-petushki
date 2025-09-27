package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/config"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/handlers"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware/auth"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware/cors"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/application"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Secret Guest API
// @version 1.0
// @description API for "Secret Guest" app
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func InitRoutes(router *gin.Engine, appUseCase application.ApplicationUseCase, cfg *config.RestConfig) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(cors.CORS(cfg.AllowOrigin))

	authProvider := auth.NewAuth()

	api := router.Group("/api/v1")

	handlers.InitApplicationHandlers(api, authProvider, appUseCase)
	handlers.InitOfferHandlers(api, authProvider)
	handlers.InitUserHandlers(api, authProvider)
	handlers.InitReportHandlers(api, authProvider)
}
