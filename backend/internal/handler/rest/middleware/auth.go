package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/dto"
)

func LoginProtected() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.TokenRequest

		if err := ctx.BindHeader(req); err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		// валидации токена

		ctx.Next()
	}
}

func RoleProtected(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.TokenRequest

		if err := ctx.BindHeader(req); err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		// валидации токена

		ctx.Next()
	}
}
