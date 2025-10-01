package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/config"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/handlers"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware/auth"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Secret Guest API
// @version 1.0
// @description API for "Secret Guest" app
// @host localhost:8081
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func initAllEndpoints(
	engine *gin.Engine,
	cfg *config.RestConfig,
	authProvider auth.Auth,
	userHandler handlers.UserHandler,
	applicationHandler handlers.ApplicationHandler,
	offerHandler handlers.OfferHandler,
	reportHandler handlers.ReportHandler,
	hotelHandler handlers.HotelHandler,
	locationHandler handlers.LocationHandler,
	roomHandler handlers.RoomHandler,
	healthHandler handlers.HealthHandler,
	client *sqlx.DB,
) {
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.Use(cors.CORS(cfg.AllowOrigin))
	engine.GET("/health", healthHandler.Health)

	router := engine.Group("/api/v1")

	initUserEndpoints(router, authProvider, userHandler)
	initApplicationHandler(router, authProvider, applicationHandler)
	initOfferHandler(router, authProvider, offerHandler)
	initReportHandler(router, authProvider, reportHandler)
	initHotelHandler(router, authProvider, hotelHandler)
	initLocationHandler(router, authProvider, locationHandler)
	initRoomHandler(router, authProvider, roomHandler)

	router.POST("test", InitDataHandler(client))
}

func initUserEndpoints(router *gin.RouterGroup, authProvider auth.Auth, h handlers.UserHandler) {
	group := router.Group("/user")

	{
		group.POST("/log-in", h.LogIn)
		group.POST("/sign-up", h.SignUp)
		group.POST("/refresh", h.Refresh)
		group.GET("/", authProvider.LoginProtected(), h.GetMe)
	}
}

func initApplicationHandler(router *gin.RouterGroup, authProvider auth.Auth, h handlers.ApplicationHandler) {

	group := router.Group("/application")

	{
		group.POST("/", authProvider.RoleProtected("reviewer"), h.CreateApplication)
		group.GET("/", authProvider.RoleProtected("reviewer"), h.GetApplications)
		group.GET("/limit", authProvider.RoleProtected("reviewer"), h.GetUserAppLimitInfo)
		group.GET("/:id", authProvider.RoleProtected("reviewer"), h.GetApplicationById)
		group.GET("/search", authProvider.RoleProtected("admin"), h.GetAppsByFilter)
	}
}

func initOfferHandler(router *gin.RouterGroup, authProvider auth.Auth, h handlers.OfferHandler) {

	group := router.Group("/offer")

	{
		group.POST("/", authProvider.RoleProtected("admin"), h.CreateOffer)
		group.GET("/", authProvider.RoleProtected("admin"), h.GetOffers)
		group.GET("/:id", authProvider.RoleProtected("admin"), h.GetOfferById)
		group.PATCH("/:id", authProvider.RoleProtected("admin"), h.UpdateOffer)

		group.GET("/search", authProvider.RoleProtected("reviewer"), h.FindOffers)
	}
}

func initReportHandler(router *gin.RouterGroup, authProvider auth.Auth, h handlers.ReportHandler) {

	group := router.Group("/report")

	{
		group.GET("/", authProvider.RoleProtected("admin"), h.GetReports)
		group.GET("/search", authProvider.RoleProtected("admin"), h.GetReportsByFilter)
		group.GET("/:id", authProvider.RoleProtected("admin"), h.GetReportById)
		group.PATCH("/:id/confirm", authProvider.RoleProtected("admin"), h.ConfirmReport)

		group.GET("/my", authProvider.RoleProtected("reviewer"), h.GetMyReports)
		group.GET("/my/:id", authProvider.RoleProtected("reviewer"), h.GetMyReportById)
		group.PATCH("/:id", authProvider.RoleProtected("reviewer"), h.UpdateReport)

		group.GET("/my/application/:id", authProvider.RoleProtected("reviewer"), h.GetMyReportByApplicationId)
	}
}

func initHotelHandler(router *gin.RouterGroup, authProvider auth.Auth, h handlers.HotelHandler) {
	group := router.Group("/hotel")

	{
		group.POST("/", authProvider.RoleProtected("admin"), h.CreateHotel)
		group.GET("/", h.GetHotels)
	}
}

func initLocationHandler(router *gin.RouterGroup, authProvider auth.Auth, h handlers.LocationHandler) {
	group := router.Group("/location")
	{
		group.POST("/", authProvider.RoleProtected("admin"), h.CreateLocation)
		group.GET("/", h.GetLocations)
	}
}

func initRoomHandler(router *gin.RouterGroup, authProvider auth.Auth, h handlers.RoomHandler) {
	group := router.Group("/room")
	{
		group.POST("/", authProvider.RoleProtected("admin"), h.CreateRoom)
		group.GET("/", h.GetRooms)
	}
}
