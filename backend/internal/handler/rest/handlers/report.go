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

	resp := &docs.GetReportsResponse{}

	ctx.JSON(http.StatusOK, resp)
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
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)

	if idStr == "" || err != nil {
		log.Println("invalid report id", idStr)
		ctx.String(http.StatusBadRequest, "invalid report id")
		return
	}

	_ = id

	resp := &docs.ReportResponse{}

	ctx.JSON(http.StatusOK, resp)
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

	userId, err := auth.GetUserId(ctx)

	if err != nil {
		log.Println("invalid user_id")
		ctx.String(http.StatusBadRequest, "invalid user_id")
		return
	}

	_ = userId

	resp := &docs.GetReportsResponse{}

	ctx.JSON(http.StatusOK, resp)
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
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)

	if idStr == "" || err != nil {
		log.Println("invalid report id", idStr)
		ctx.String(http.StatusBadRequest, "invalid report id")
		return
	}

	_ = id

	userId, err := auth.GetUserId(ctx)

	if err != nil {
		log.Println("invalid user_id")
		ctx.String(http.StatusBadRequest, "invalid user_id")
		return
	}

	_ = userId

	resp := &docs.GetReportsResponse{}

	ctx.JSON(http.StatusOK, resp)
}

// Add godoc
// @Summary Update report
// @Description Updates report with given text and photos
// @Tags Report
// @Accept multipart/form-data
// @Param id path string true "Id of report to update"
// @Param text formData string true "Report text"
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
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)

	if idStr == "" || err != nil {
		log.Println("invalid report id", idStr)
		ctx.String(http.StatusBadRequest, "invalid report id")
		return
	}

	_ = id

	userId, err := auth.GetUserId(ctx)

	if err != nil {
		log.Println("invalid user_id")
		ctx.String(http.StatusBadRequest, "invalid user_id")
		return
	}

	_ = userId

	var request docs.UpdateReportRequest

	if err := ctx.ShouldBind(&request); err != nil {
		log.Println("invalid form data")
		ctx.String(http.StatusBadRequest, "invalud form data")
		return
	}

	ctx.Status(http.StatusOK)
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
	var request docs.ConfirmReport

	if err := ctx.BindJSON(&request); err != nil {
		log.Println("Invalid body")
		ctx.String(http.StatusBadRequest, "invalid body")
		return
	}

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)

	if idStr == "" || err != nil {
		log.Println("invalid report id", idStr)
		ctx.String(http.StatusBadRequest, "invalid report id")
		return
	}

	_ = id

	ctx.Status(http.StatusOK)
}

func InitReportHandlers(router *gin.RouterGroup, authProvider *auth.Auth) {
	h := &ReportHandlers{}

	group := router.Group("/report")

	{
		group.GET("/", authProvider.RoleProtected("admin"), h.GetReports)
		group.GET("/:id", authProvider.RoleProtected("admin"), h.GetReportById)
		group.PATCH("/:id/confirm", authProvider.RoleProtected("admin"), h.ConfirmReport)

		group.GET("/my", authProvider.RoleProtected("reviewer"), h.GetMyReports)
		group.GET("/my/:id", authProvider.RoleProtected("reviewer"), h.GetMyReportById)
		group.PATCH("/:id", authProvider.RoleProtected("reviewer"), h.UpdateReport)
	}
}
