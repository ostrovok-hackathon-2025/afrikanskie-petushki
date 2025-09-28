package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware/auth"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/validation"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/user"
)

type UserHandler interface {
	LogIn(ctx *gin.Context)
	SignUp(ctx *gin.Context)
	Refresh(ctx *gin.Context)
	GetMe(ctx *gin.Context)
}

type userHandler struct {
	useCase user.UseCase
}

func NewUserHandler(useCase user.UseCase) UserHandler {
	return &userHandler{
		useCase: useCase,
	}
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
func (h *userHandler) LogIn(ginCtx *gin.Context) {
	var request docs.LogInRequest
	ctx := context.Background()

	if err := ginCtx.BindJSON(&request); err != nil {
		log.Println("Invalid body")
		ginCtx.String(http.StatusBadRequest, "invalid body")
		return
	}

	if err := validation.ValidateUsername(request.OstrovokLogin); err != nil {
		log.Println("Invalid username: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := validation.ValidatePassword(request.Password); err != nil {
		log.Println("Invalid username: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.useCase.Login(ctx, &request)

	if err != nil {
		//TODO обработка похитрее
		ginCtx.JSON(http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, resp)
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
func (h *userHandler) SignUp(ginCtx *gin.Context) {
	var request docs.SignUpRequest
	ctx := context.Background()

	if err := ginCtx.BindJSON(&request); err != nil {
		log.Println("Invalid body")
		ginCtx.String(http.StatusBadRequest, "invalid body")
		return
	}

	if err := validation.ValidateUsername(request.OstrovokLogin); err != nil {
		log.Println("Invalid username: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := validation.ValidateEmail(request.Email); err != nil {
		log.Println("Invalid email: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := validation.ValidatePassword(request.Password); err != nil {
		log.Println("Invalid password: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.useCase.Register(ctx, &request)

	if err != nil {
		//TODO обработка похитрее
		log.Println("Err from useCAse password: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}

	ginCtx.JSON(http.StatusCreated, resp)
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
func (h *userHandler) Refresh(ctx *gin.Context) {
	var request docs.RefreshRequest

	if err := ctx.BindJSON(&request); err != nil {
		log.Println("Invalid body")
		ctx.String(http.StatusBadRequest, "invalid body")
		return
	}

	resp := &docs.AuthResponse{}

	ctx.JSON(http.StatusOK, resp)
}

// Add godoc
// @Summary GetForPage me
// @Description GetForPage data of current user
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.UserResponse "User data"
// @Failure 400 {string} string "Bad request to get user data"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /user/ [get]
func (h *userHandler) GetMe(ctx *gin.Context) {
	userId, err := auth.GetUserId(ctx)

	if err != nil {
		log.Println("invalid user_id")
		ctx.String(http.StatusBadRequest, "invalid user_id")
		return
	}

	_ = userId

	resp := &docs.UserResponse{}

	ctx.JSON(http.StatusOK, resp)
}
