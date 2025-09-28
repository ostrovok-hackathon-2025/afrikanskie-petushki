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

type ReportHandler interface {
	GetReports(ctx *gin.Context)
	GetReportById(ctx *gin.Context)
	GetMyReports(ctx *gin.Context)
	GetMyReportById(ctx *gin.Context)
	UpdateReport(ctx *gin.Context)
	ConfirmReport(ctx *gin.Context)
}

type reportHandler struct {
}

func NewReportHandler() ReportHandler {
	return &reportHandler{}
}

// Add godoc
// @Summary GetForPage reports
// @Description GetForPage all reports with pagination
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
func (h *reportHandler) GetReports(ctx *gin.Context) {
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
// @Summary GetForPage by id
// @Description GetForPage report by id
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
func (h *reportHandler) GetReportById(ctx *gin.Context) {
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
// @Summary GetForPage my reports
// @Description GetForPage all reports of current user with pagination
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
func (h *reportHandler) GetMyReports(ctx *gin.Context) {
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
// @Summary GetForPage my by id
// @Description GetForPage report of current user by id
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
func (h *reportHandler) GetMyReportById(ctx *gin.Context) {
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
func (h *reportHandler) UpdateReport(ctx *gin.Context) {
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
// @Param id path string true "Id of report to update"
// @Param input body docs.ConfirmReport true "Data for updating report"
// @Produce json
// @Security BearerAuth
// @Success 200 "Successfully update status of report"
// @Failure 400 {string} string "Invalid data for changing report status"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 500 "Internal server error"
// @Router /report/{id}/confirm [patch]
func (h *reportHandler) ConfirmReport(ctx *gin.Context) {
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
