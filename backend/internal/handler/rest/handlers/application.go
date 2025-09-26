package handlers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware"
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
// @Failure 403 "Only available for reviewer"
// @Failure 404 "Application with given id not found"
// @Failure 500 "Internal server error"
// @Router /application/{id} [get]
func (h *ApplicationHandlers) GetApplicationById(ctx *gin.Context) {

}

func InitApplicationHandlers(router *gin.RouterGroup) {
	h := &ApplicationHandlers{}

	group := router.Group("/application")

	group.Use(middleware.RoleProtected("reviewer"))

	{
		group.POST("/", h.CreateApplication)
		group.GET("/", h.GetApplications)
		group.GET("/:id", h.GetApplicationById)
	}
}
