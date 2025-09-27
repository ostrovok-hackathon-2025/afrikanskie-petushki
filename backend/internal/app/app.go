package app

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	applicationRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/application"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/offer"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/config"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/application"
)

func MustConfigureApp(router *gin.Engine, cfg *config.Config) func() {

	logger := log.New(os.Stdout, cfg.LoggerConfig.Prefix, cfg.LoggerConfig.Flag)
	postgresClient := initPostgresClient(&cfg.PostgresConfig)
	offer.New(postgresClient, logger)
	var appRepo applicationRepo.ApplicationRepo = nil
	appUseCase := application.NewApplicationService(appRepo)

	rest.InitRoutes(router, appUseCase, &cfg.RestConfig)

	return func() {}
}
