package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware"
)

type UserHandlers struct {
}

// Add godoc
// @Summary Log in
// @Description Log in with given name and password
// @Tags User
// @Accept json
// @Param input body docs.LogInRequest true "Data for log in"
// @Produce json
// @Success 200 {object} docs.AuthResponse "Auth data"
// @Failure 400 {string} string "Invalid data for log in"
// @Failure 401 "Wrong password"
// @Failure 404 "User not found"
// @Failure 500 "Internal server error"
// @Router /user/log-in [post]
func (h *UserHandlers) LogIn(ctx *gin.Context) {

}

// Add godoc
// @Summary Sign up
// @Description Sign up with given name and password
// @Tags User
// @Accept json
// @Param input body docs.SignUpRequest true "Data for sign up"
// @Produce json
// @Success 201 {object} docs.AuthResponse "Auth data"
// @Failure 400 {string} string "Invalid data for sign up"
// @Failure 404 "User not found"
// @Failure 500 "Internal server error"
// @Router /user/sign-up [post]
func (h *UserHandlers) SignUp(ctx *gin.Context) {

}

// Add godoc
// @Summary Refresh credentials
// @Description Refresh auth credentials via refresh token
// @Tags User
// @Accept json
// @Param input body docs.RefreshRequest true "Data for refresh"
// @Produce json
// @Success 200 {object} docs.AuthResponse "Auth data"
// @Failure 400 {string} string "Invalid data for refresh"
// @Failure 500 "Internal server error"
// @Router /user/refresh [post]
func (h *UserHandlers) Refresh(ctx *gin.Context) {

}

// Add godoc
// @Summary Get me
// @Description Get data of current user
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.UserResponse "User data"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /user/ [get]
func (h *UserHandlers) GetMe(ctx *gin.Context) {

}

func InitUserHandlers(router *gin.RouterGroup) {
	h := &UserHandlers{}

	group := router.Group("/user")

	{
		group.POST("/log-in", h.LogIn)
		group.POST("/sign-up", h.SignUp)
		group.POST("/refresh", h.Refresh)
		group.GET("/", middleware.LoginProtected(), h.GetMe)
	}
}
