package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware"
)

type ReportHandlers struct {
}

// Add godoc
// @Summary Get reports
// @Description Get all reports with pagination
// @Tags Report
// @Param pageNum query int true "Number of page"
// @Param pageSize query int true "Size of page"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.GetReportsResponse "Page of reports"
// @Failure 400 {string} string "Invalid data for getting reports"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 404 "Page with given number not found"
// @Failure 500 "Internal server error"
// @Router /report/ [get]
func (h *ReportHandlers) GetReports(ctx *gin.Context) {

}

// Add godoc
// @Summary Get by id
// @Description Get report by id
// @Tags Report
// @Param id path string true "Id of requested report"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.ReportResponse "Requested report"
// @Failure 400 {string} string "Invalid data for getting report by id"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 404 "Report with given id not found"
// @Failure 500 "Internal server error"
// @Router /report/{id} [get]
func (h *ReportHandlers) GetReportById(ctx *gin.Context) {

}

// Add godoc
// @Summary Get my reports
// @Description Get all reports of current user with pagination
// @Tags Report
// @Param pageNum query int true "Number of page"
// @Param pageSize query int true "Size of page"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.GetReportsResponse "Page of reports"
// @Failure 400 {string} string "Invalid data for getting reports"
// @Failure 401 "Unauthorized"
// @Failure 403 "User is not reviewer"
// @Failure 404 "Page with given number not found"
// @Failure 500 "Internal server error"
// @Router /report/my [get]
func (h *ReportHandlers) GetMyReports(ctx *gin.Context) {

}

// Add godoc
// @Summary Get my by id
// @Description Get report of current user by id
// @Tags Report
// @Param id path string true "Id of requested report"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.ReportResponse "Requested report"
// @Failure 400 {string} string "Invalid data for getting report by id"
// @Failure 401 "Unauthorized"
// @Failure 403 {string} string "User is not reviewer or this report does not belong to user"
// @Failure 404 "Report with given id not found"
// @Failure 500 "Internal server error"
// @Router /report/my/{id} [get]
func (h *ReportHandlers) GetMyReportById(ctx *gin.Context) {

}

// Add godoc
// @Summary Update report
// @Description Updates report with given text and photos
// @Tags Report
// @Accept multipart/form-data
// @Param id path string true "Id of report to update"
// @Param input body docs.UpdateReportRequest true "Data for updating report"
// @Param images formData file true "Report images" collectionFormat multi
// @Produce json
// @Security BearerAuth
// @Success 200 "Successfully update report"
// @Failure 400 {string} string "Invalid data for updating report"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for reviewer"
// @Failure 500 "Internal server error"
// @Router /report/{id} [patch]
func (h *ReportHandlers) UpdateReport(ctx *gin.Context) {

}

// Add godoc
// @Summary Confirm report
// @Description Confirms or declains report
// @Tags Report
// @Accept json
// @Param input body docs.ConfirmReport true "Data for updating report"
// @Produce json
// @Security BearerAuth
// @Success 200 "Successfully update status of report"
// @Failure 400 {string} string "Invalid data for changing report status"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 500 "Internal server error"
// @Router /report/{id}/confirm [patch]
func (h *ReportHandlers) ConfirmReport(ctx *gin.Context) {

}

func InitReportHandlers(router *gin.RouterGroup) {
	h := &ReportHandlers{}

	group := router.Group("/report")

	{
		group.GET("/", middleware.RoleProtected("admin"), h.GetReports)
		group.GET("/:id", middleware.RoleProtected("admin"), h.GetReportById)
		group.PATCH("/:id/confirm", middleware.RoleProtected("admin"), h.ConfirmReport)

		group.GET("/my", middleware.RoleProtected("reviewer"), h.GetMyReports)
		group.GET("/my/:id", middleware.RoleProtected("reviewer"), h.GetMyReportById)
		group.PATCH("/:id", middleware.RoleProtected("reviewer"), h.UpdateReport)
	}
}
