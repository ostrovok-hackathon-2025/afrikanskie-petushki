package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/config"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest"
)

func MustConfigureApp(router *gin.Engine, cfg *config.Config) func() {
	rest.InitRoutes(router, &cfg.RestConfig)

	return func() {}
}
