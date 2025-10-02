package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type HealthHandler interface {
	Health(ctx *gin.Context)
}

type healthHandler struct {
	sqlClient   *sqlx.DB
	minioClient *minio.Client
}

func NewHealthHandler(sqlClient *sqlx.DB, minioClient *minio.Client) HealthHandler {
	return &healthHandler{
		sqlClient:   sqlClient,
		minioClient: minioClient,
	}
}

func (h *healthHandler) Health(ctx *gin.Context) {
	err := h.sqlClient.Ping()
	if err != nil {
		log.Println("Err to ping sqlClient: ", err.Error())
		ctx.String(http.StatusBadRequest, "sql client not ready")
	}
	_, err = h.minioClient.ListBuckets(ctx)
	if err != nil {
		log.Println("Err to get ListBuckets: ", err.Error())
		ctx.String(http.StatusBadRequest, "minio client not ready")
	}
	ctx.String(http.StatusOK, "ok")
}
