package app

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/ostrovok"
	applicationRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/application"
	offerRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/offer"
	userRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/user"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/config"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/handlers"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware/auth"
	applicationUC "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/application"
	offerUC "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/offer"
	userUC "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/user"
)

func MustConfigureApp(engine *gin.Engine, cfg *config.Config) func() {

	logger := log.New(os.Stdout, cfg.LoggerConfig.Prefix, cfg.LoggerConfig.Flag)

	//Clients
	ostrovokClient := ostrovok.NewClient()
	sqlClient := initPostgresClient(&cfg.PostgresConfig)

	//Repos

	applicationRepository := applicationRepo.NewApplicationRepo(sqlClient)
	userRepository := userRepo.NewUserRepo(sqlClient)
	offerRepository := offerRepo.New(sqlClient, logger)

	//UseCases

	applicationService := applicationUC.NewApplicationService(applicationRepository)
	userUseCase := userUC.NewUseCase(userRepository, ostrovokClient)
	offerUseCase := offerUC.NewUseCase(offerRepository)

	//Handlers

	userHandler := handlers.NewUserHandler(userUseCase)
	applicationHandler := handlers.NewApplicationHandler(applicationService)
	offerHandler := handlers.NewOfferHandler(offerUseCase)
	reportHandler := handlers.NewReportHandler()

	//MiddleWare
	authMiddleWare := auth.NewAuth()

	//InitEndpoints
	initAllEndpoints(engine, &cfg.RestConfig, authMiddleWare, userHandler, applicationHandler, offerHandler, reportHandler)

	return func() {}
}
