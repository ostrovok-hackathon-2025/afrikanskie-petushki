package app

import (
	"log"
	"os"

	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/ostrovok"
	analyticsRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/analytics"
	applicationRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/application"
	hotelRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/hotel"
	locationRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/location"
	offerRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/offer"
	reportRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/report"
	roomRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/room"
	userRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/user"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/s3/image"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/config"

	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/handlers"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware/auth"
	analyticsUC "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/analytics"
	applicationUC "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/application"
	hotelUC "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/hotel"
	locationUC "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/location"
	offerUC "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/offer"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/report"
	roomUC "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/room"
	userUC "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/user"
)

func MustConfigureApp(engine *gin.Engine, cfg *config.Config) func() {

	logger := log.New(os.Stdout, cfg.LoggerConfig.Prefix, cfg.LoggerConfig.Flag)

	//Clients
	ostrovokClient := ostrovok.NewClient()
	sqlClient := initPostgresClient(&cfg.PostgresConfig)
	minioClient, err := initMinioConnection(&cfg.MinioConfig)

	if err != nil {
		log.Fatalf("failed to connect to minio: %s", err.Error())
	}

	//Repos

	applicationRepository := applicationRepo.NewApplicationRepo(sqlClient)
	userRepository := userRepo.NewUserRepo(sqlClient)
	offerRepository := offerRepo.New(sqlClient, logger)
	hotelRepository := hotelRepo.NewRepo(sqlClient)
	locationRepository := locationRepo.NewRepo(sqlClient)
	roomRepository := roomRepo.NewRepo(sqlClient)
	reportRepository := reportRepo.NewRepo(sqlClient)
	analyticsRepository := analyticsRepo.NewRepo(sqlClient)

	imageRepo := image.NewImageRepoMinio(minioClient, cfg.MinioConfig.PublicEndpoint, cfg.MinioConfig.BucketName)

	//UseCases

	applicationService := applicationUC.NewApplicationService(applicationRepository)
	userUseCase := userUC.NewUseCase(userRepository, ostrovokClient)
	offerUseCase := offerUC.NewUseCase(offerRepository)
	hotelUseCase := hotelUC.NewUseCase(hotelRepository)
	locationUseCase := locationUC.NewUseCase(locationRepository)
	roomUseCase := roomUC.NewUseCase(roomRepository)
	reportUsccase := report.New(reportRepository, imageRepo, ostrovokClient, userRepository, applicationRepository)
	analyticsUseCase := analyticsUC.NewAnalyticsUseCase(analyticsRepository)

	//Handlers

	userHandler := handlers.NewUserHandler(userUseCase)
	applicationHandler := handlers.NewApplicationHandler(applicationService)
	offerHandler := handlers.NewOfferHandler(offerUseCase)
	reportHandler := handlers.NewReportHandler(reportUsccase)
	hotelHandler := handlers.NewHotelHandler(hotelUseCase)
	locationHandler := handlers.NewLocationHandler(locationUseCase)
	roomHandler := handlers.NewRoomHandler(roomUseCase)
	heathHandler := handlers.NewHealthHandler(sqlClient, minioClient)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsUseCase)

	//MiddleWare
	authMiddleWare := auth.NewAuth(userUseCase)

	//InitEndpoints
	initAllEndpoints(
		engine,
		&cfg.RestConfig,
		authMiddleWare,
		userHandler,
		applicationHandler,
		offerHandler,
		reportHandler,
		hotelHandler,
		locationHandler,
		roomHandler,
		analyticsHandler,
		heathHandler,
		sqlClient,
	)

	secretGuestWorker := worker.NewSecretGuestWorker(offerRepository, applicationRepository, reportRepository)
	secretGuestWorker.Start()

	return func() {
		secretGuestWorker.Stop()
	}
}
