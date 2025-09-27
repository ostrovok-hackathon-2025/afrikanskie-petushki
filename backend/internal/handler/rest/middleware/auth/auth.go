package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware/auth/dto"
)

const _AUTH_USER_ID = "__auth_user_id"
const _AUTH_USER_ROLE = "__auth_user_role"

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) LoginProtected() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.TokenRequest

		if err := ctx.BindHeader(req); err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userId, userRole, err := a.parseToken(req.Auth)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set(_AUTH_USER_ID, userId)
		ctx.Set(_AUTH_USER_ROLE, userRole)

		ctx.Next()
	}
}

func (a *Auth) RoleProtected(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.TokenRequest

		if err := ctx.BindHeader(req); err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		userId, userRole, err := a.parseToken(req.Auth)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if userRole != role {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Set(_AUTH_USER_ID, userId)
		ctx.Set(_AUTH_USER_ROLE, userRole)

		ctx.Next()
	}
}

func (a *Auth) parseToken(token string) (uuid.UUID, string, error) {
	// валидации токена
	return uuid.UUID{}, "", nil
}

func GetUserId(ctx *gin.Context) (uuid.UUID, error) {
	idStr := ctx.GetString(_AUTH_USER_ID)
	id, err := uuid.Parse(idStr)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to parse user uuid: %w", err)
	}

	return id, nil
}

func GetUserRole(ctx *gin.Context) (string, error) {
	role := ctx.GetString(_AUTH_USER_ROLE)

	return role, nil
}
