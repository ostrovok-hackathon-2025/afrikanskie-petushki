package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware/auth/dto"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/user"
)

const _AUTH_USER_ID = "__auth_user_id"
const _AUTH_USER_ROLE = "__auth_user_role"

type Auth interface {
	LoginProtected() gin.HandlerFunc
	RoleProtected(role string) gin.HandlerFunc
	parseToken(token string) (string, string, error)
}

type auth struct {
	uc user.UseCase
}

func NewAuth(uc user.UseCase) Auth {
	return &auth{
		uc: uc,
	}
}

func (a *auth) LoginProtected() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.TokenRequest

		if err := ctx.BindHeader(&req); err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userId, userRole, err := a.parseToken(req.Auth)
		if err != nil {
			log.Println("failed to read token", err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set(_AUTH_USER_ID, userId)
		ctx.Set(_AUTH_USER_ROLE, userRole)

		ctx.Next()
	}
}

func (a *auth) RoleProtected(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.TokenRequest

		if err := ctx.BindHeader(&req); err != nil {
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

func (a *auth) parseToken(tokenQuery string) (string, string, error) {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(tokenQuery, bearerPrefix) {
		return "", "", errors.New("invalid token")
	}

	token := strings.TrimPrefix(tokenQuery, bearerPrefix)
	if token == "" {
		return "", "", errors.New("invalid token")
	}

	claims, err := a.uc.ValidateToken(token)

	if err != nil {
		return "", "", fmt.Errorf("failed to validate token: %w", err)
	}

	role := "reviewer"
	if claims.IsAdmin {
		role = "admin"
	}

	return claims.UserID, role, nil
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
