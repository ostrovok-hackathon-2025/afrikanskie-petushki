package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware/auth"
)

type ApplicationHandlers struct {
}

// Add godoc
// @Summary Create application
// @Description Creates application for given offer
// @Tags Application
// @Accept json
// @Param input body docs.CreateApplicationRequest true "Data for creating offer"
// @Produce json
// @Security BearerAuth
// @Success 201 {object} docs.CreateApplicationResponse "Created application data"
// @Failure 400 {string} string "Invalid data for creating offer"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for reviewer"
// @Failure 500 "Internal server error"
// @Router /application/ [post]
func (h *ApplicationHandlers) CreateApplication(ctx *gin.Context) {
	var request docs.CreateApplicationRequest

	if err := ctx.BindJSON(&request); err != nil {
		log.Println("Invalid body")
		ctx.String(http.StatusBadRequest, "invalid body")
		return
	}

	offerIdStr := request.OfferId
	offerId, err := uuid.Parse(offerIdStr)

	if offerIdStr == "" || err != nil {
		log.Println("Invalid offer_id: ", offerIdStr)
		ctx.String(http.StatusBadRequest, "invalid offer_id")
		return
	}

	_ = offerId

	userId, err := auth.GetUserId(ctx)

	if err != nil {
		log.Println("invalid user_id")
		ctx.String(http.StatusBadRequest, "invalid user_id")
	}

	_ = userId

	resp := &docs.CreateApplicationResponse{}

	ctx.JSON(http.StatusCreated, resp)
}

// Add godoc
// @Summary Get applications
// @Description Get all applications with pagination
// @Tags Application
// @Param pageNum query int true "Number of page"
// @Param pageSize query int true "Size of page"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.GetApplicationsResponse "Page of applications"
// @Failure 400 {string} string "Invalid data for getting applications"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for reviewer"
// @Failure 404 "Page with given number not found"
// @Failure 500 "Internal server error"
// @Router /application/ [get]
func (h *ApplicationHandlers) GetApplications(ctx *gin.Context) {
	pageNumStr := ctx.Query("pageNum")
	pageNum, err := strconv.Atoi(pageNumStr)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageNum: ", pageNumStr)
		ctx.String(http.StatusBadRequest, "invalid pageNum")
		return
	}

	_ = pageNum

	pageSizeStr := ctx.Query("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageSize: ", pageSizeStr)
		ctx.String(http.StatusBadRequest, "invalid pageSize")
		return
	}

	_ = pageSize

	resp := &docs.GetApplicationsResponse{}

	ctx.JSON(http.StatusOK, resp)
}

// Add godoc
// @Summary Get by id
// @Description Get application by id
// @Tags Application
// @Param id path string true "Id of requested application"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.ApplicationResponse "Requested application"
// @Failure 400 {string} string "Invalid data for getting application by id"
// @Failure 401 "Unauthorized"
// @Failure 403 "User is not reviewer or application does not belong to user"
// @Failure 404 "Application with given id not found"
// @Failure 500 "Internal server error"
// @Router /application/{id} [get]
func (h *ApplicationHandlers) GetApplicationById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)

	if idStr == "" || err != nil {
		log.Println("invalid application id", idStr)
		ctx.String(http.StatusBadRequest, "invalid application id")
		return
	}

	_ = id

	resp := &docs.ApplicationResponse{}

	ctx.JSON(http.StatusOK, resp)
}

func InitApplicationHandlers(router *gin.RouterGroup, authProvider *auth.Auth) {
	h := &ApplicationHandlers{}

	group := router.Group("/application")

	group.Use(authProvider.RoleProtected("reviewer"))

	{
		group.POST("/", h.CreateApplication)
		group.GET("/", h.GetApplications)
		group.GET("/:id", h.GetApplicationById)
	}
}
