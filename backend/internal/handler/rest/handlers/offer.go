package handlers

import (
	"github.com/gin-gonic/gin"

	_ "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/middleware"
)

type OfferHandlers struct {
}

// Add godoc
// @Summary Create offer
// @Description Creates offer with given info
// @Tags Offer
// @Accept json
// @Param input body docs.CreateOfferRequest true "Data for creating offer"
// @Produce json
// @Security BearerAuth
// @Success 201 {object} docs.CreateOfferResponse "Created offer data"
// @Failure 400 {string} string "Invalid data for creating offer"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 500 "Internal server error"
// @Router /offer/ [post]
func (h *OfferHandlers) CreateOffer(ctx *gin.Context) {

}

// Add godoc
// @Summary Get offers
// @Description Get all offers with pagination
// @Tags Offer
// @Param pageNum query int true "Number of page"
// @Param pageSize query int true "Size of page"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.GetOffersResponse "Page of offers"
// @Failure 400 {string} string "Invalid data for getting offers"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 404 "Page with given number not found"
// @Failure 500 "Internal server error"
// @Router /offer/ [get]
func (h *OfferHandlers) GetOffers(ctx *gin.Context) {

}

// Add godoc
// @Summary Get by id
// @Description Get offer by id
// @Tags Offer
// @Param id path string true "Id of requested offer"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.OfferResponse "Requested offer"
// @Failure 400 {string} string "Invalid data for getting offer by id"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 404 "Offer with given id not found"
// @Failure 500 "Internal server error"
// @Router /offer/{id} [get]
func (h *OfferHandlers) GetOfferById(ctx *gin.Context) {

}

// Add godoc
// @Summary Find offers
// @Description Find offers with given search params
// @Tags Offer
// @Param cityId query string true "Id of required city"
// @Param pageNum query int true "Number of page"
// @Param pageSize query int true "Size of page"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.GetOffersResponse "Page of offers found"
// @Failure 400 {string} string "Invalid data for finding offers"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for reviewer"
// @Failure 404 "Page with given number not found"
// @Failure 500 "Internal server error"
// @Router /offer/search [get]
func (h *OfferHandlers) FindOffers(ctx *gin.Context) {

}

// Add godoc
// @Summary Update offer
// @Description Update offer with given id and data
// @Tags Offer
// @Accept json
// @Param id path string true "Id of offer to update"
// @Param input body docs.UpdateOfferRequest true "Data for updating offer"
// @Produce json
// @Security BearerAuth
// @Success 200 "Successfully updated offer"
// @Failure 400 {string} string "Invalid data for finding offers"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 404 "Offer with given id not found"
// @Failure 500 "Internal server error"
// @Router /offer/{id} [patch]
func (h *OfferHandlers) UpdateOffer(ctx *gin.Context) {

}

func InitOfferHandlers(router *gin.RouterGroup) {
	h := &OfferHandlers{}

	group := router.Group("/offer")

	{
		group.POST("/", middleware.RoleProtected("admin"), h.CreateOffer)
		group.GET("/", middleware.RoleProtected("admin"), h.GetOffers)
		group.GET("/:id", middleware.RoleProtected("admin"), h.GetOfferById)
		group.PATCH("/:id", middleware.RoleProtected("admin"), h.UpdateOffer)

		group.GET("/search", middleware.RoleProtected("reviewer"), h.FindOffers)
	}
}
