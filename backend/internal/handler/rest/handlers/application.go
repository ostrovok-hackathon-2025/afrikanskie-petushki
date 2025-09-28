package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware/auth"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/application"

	applicationRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/application"
)

type ApplicationHandler interface {
	CreateApplication(ctx *gin.Context)
	GetApplications(ctx *gin.Context)
	GetApplicationById(ctx *gin.Context)
}

type applicationHandler struct {
	useCase application.ApplicationUseCase
}

func NewApplicationHandler(useCase application.ApplicationUseCase) ApplicationHandler {
	return &applicationHandler{
		useCase: useCase,
	}
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
func (h *applicationHandler) CreateApplication(ctx *gin.Context) {
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

	userId, err := auth.GetUserId(ctx)

	if err != nil {
		log.Println("invalid user_id")
		ctx.String(http.StatusBadRequest, "invalid user_id")
	}

	appId, err := h.useCase.CreateApplication(ctx.Request.Context(), userId, offerId)

	if err != nil {
		log.Println("failed to create app", err)

		switch {
		case errors.Is(err, applicationRepo.ErrOfferNotExist):
			ctx.String(http.StatusBadRequest, "offer does not exist")
		case errors.Is(err, applicationRepo.ErrUserNotExist):
			ctx.String(http.StatusBadRequest, "user does not exist")
		default:
			ctx.Status(http.StatusInternalServerError)
		}

		return
	}

	resp := &docs.CreateApplicationResponse{
		ApplicationId: appId.String(),
	}

	ctx.JSON(http.StatusCreated, resp)
}

// Add godoc
// @Summary GetForPage applications
// @Description GetForPage all applications with pagination
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
func (h *applicationHandler) GetApplications(ctx *gin.Context) {
	pageNumStr := ctx.Query("pageNum")
	pageNum, err := strconv.Atoi(pageNumStr)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageNum: ", pageNumStr)
		ctx.String(http.StatusBadRequest, "invalid pageNum")
		return
	}

	pageSizeStr := ctx.Query("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageSize: ", pageSizeStr)
		ctx.String(http.StatusBadRequest, "invalid pageSize")
		return
	}

	userId, err := auth.GetUserId(ctx)

	if err != nil {
		log.Println("invalid user_id")
		ctx.String(http.StatusBadRequest, "invalid user_id")
		return
	}

	apps, count, err := h.useCase.GetApplications(ctx.Request.Context(), userId, pageNum, pageSize)

	if err != nil {
		log.Println("failed to get all applications", err)

		switch {
		case errors.Is(err, applicationRepo.ErrPageNotFound):
			ctx.Status(http.StatusNotFound)
		default:
			ctx.Status(http.StatusInternalServerError)
		}

		return
	}

	appsResp := make([]*docs.ApplicationResponse, 0, len(apps))

	for _, app := range apps {
		appsResp = append(appsResp, docs.ApplicationModelToResponse(app))
	}

	resp := &docs.GetApplicationsResponse{
		Applications: appsResp,
		PagesCount:   count,
	}

	ctx.JSON(http.StatusOK, resp)
}

// Add godoc
// @Summary GetForPage by id
// @Description GetForPage application by id
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
func (h *applicationHandler) GetApplicationById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)

	if idStr == "" || err != nil {
		log.Println("invalid application id", idStr)
		ctx.String(http.StatusBadRequest, "invalid application id")
		return
	}

	userId, err := auth.GetUserId(ctx)

	if err != nil {
		log.Println("invalid user_id")
		ctx.String(http.StatusBadRequest, "invalid user_id")
		return
	}

	app, err := h.useCase.GetApplicationById(ctx, userId, id)

	if err != nil {
		log.Println("failed to get application by id", err)

		switch {
		case errors.Is(err, applicationRepo.ErrApplicationNotFound):
			ctx.Status(http.StatusNotFound)
		case errors.Is(err, application.ErrNotOwner):
			ctx.Status(http.StatusForbidden)
		default:
			ctx.Status(http.StatusInternalServerError)
		}

		return
	}

	resp := docs.ApplicationModelToResponse(app)

	ctx.JSON(http.StatusOK, resp)
}
