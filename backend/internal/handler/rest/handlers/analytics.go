package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/analytics"
)

type AnalyticsHandler interface {
	GetAnalytics(ctx *gin.Context)
}

type analyticsHandler struct {
	useCase analytics.UseCase
}

func NewAnalyticsHandler(useCase analytics.UseCase) AnalyticsHandler {
	return &analyticsHandler{
		useCase: useCase,
	}
}

// Add godoc
// @Summary Get analytics
// @Description Calc and return metrics for analytics
// @Tags Analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.AnalyticsResponse "Page of applications"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 500 "Internal server error"
// @Router /analytics/ [get]
func (h *analyticsHandler) GetAnalytics(ctx *gin.Context) {
	analyticsRes, err := h.useCase.GetAnalytics(ctx.Request.Context())

	if err != nil {
		fmt.Println("failed to get analytics", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	resp := &docs.AnalyticsResponse{
		CompletedOffers:      analyticsRes.CompletedOffers,
		ApplicationsReceived: analyticsRes.ApplicationsReceived,
		AcceptedReports:      analyticsRes.AcceptedReports,
	}

	ctx.JSON(http.StatusOK, resp)
}
