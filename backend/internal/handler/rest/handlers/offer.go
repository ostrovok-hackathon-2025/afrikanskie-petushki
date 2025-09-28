package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/offer"

	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
)

type OfferHandler interface {
	CreateOffer(ctx *gin.Context)
	GetOffers(ctx *gin.Context)
	GetOfferById(ctx *gin.Context)
	FindOffers(ctx *gin.Context)
	UpdateOffer(ctx *gin.Context)
}

type offerHandler struct {
	useCase offer.UseCase
}

func NewOfferHandler(useCase offer.UseCase) OfferHandler {
	return &offerHandler{
		useCase: useCase,
	}
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
func (h *offerHandler) CreateOffer(ginCtx *gin.Context) {
	var request docs.CreateOfferRequest
	ctx := context.Background()

	if err := ginCtx.BindJSON(&request); err != nil {
		log.Println("Invalid body")
		ginCtx.String(http.StatusBadRequest, "invalid body")
		return
	}

	create := model.Create{}

	id, err := h.useCase.Create(ctx, &create)

	if err != nil {
		log.Println("Err to create offer: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := &docs.CreateOfferResponse{
		Id: id.String(),
	}

	ginCtx.JSON(http.StatusCreated, resp)
}

// Add godoc
// @Summary GetForPage offers
// @Description GetForPage all offers with pagination
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
func (h *offerHandler) GetOffers(ginCtx *gin.Context) {
	ctx := context.Background()
	pageNumStr := ginCtx.Query("pageNum")
	pageNum, err := strconv.Atoi(pageNumStr)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageNum: ", pageNumStr)
		ginCtx.String(http.StatusBadRequest, "invalid pageNum")
		return
	}

	pageSizeStr := ginCtx.Query("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageSize: ", pageSizeStr)
		ginCtx.String(http.StatusBadRequest, "invalid pageSize")
		return
	}
	pageSettings := model.PageSettings{
		Limit:  pageSize,
		Offset: (pageNum - 1) * pageSize,
	}
	//TODO CREATE BODY
	_, _, err = h.useCase.GetForPage(ctx, &pageSettings)

	if err != nil {
		log.Println("Err to get offers for page: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := &docs.GetOffersResponse{}

	ginCtx.JSON(http.StatusOK, resp)
}

// Add godoc
// @Summary GetForPage by id
// @Description GetForPage offer by id
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
func (h *offerHandler) GetOfferById(ginCtx *gin.Context) {
	ctx := context.Background()
	idStr := ginCtx.Param("id")
	id, err := uuid.Parse(idStr)

	if idStr == "" || err != nil {
		log.Println("invalid offer id", idStr)
		ginCtx.String(http.StatusBadRequest, "invalid offer id")
		return
	}

	//TODO CREATE BODY
	_, err = h.useCase.GetByID(ctx, id.String())
	if err != nil {
		log.Println("Err to get offer by id: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := &docs.OfferResponse{}

	ginCtx.JSON(http.StatusOK, resp)
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
func (h *offerHandler) FindOffers(ginCtx *gin.Context) {
	ctx := context.Background()
	pageNumStr := ginCtx.Query("pageNum")
	pageNum, err := strconv.Atoi(pageNumStr)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageNum: ", pageNumStr)
		ginCtx.String(http.StatusBadRequest, "invalid pageNum")
		return
	}

	_ = pageNum

	pageSizeStr := ginCtx.Query("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageSize: ", pageSizeStr)
		ginCtx.String(http.StatusBadRequest, "invalid pageSize")
		return
	}

	_ = pageSize

	cityIdStr := ginCtx.Query("pageSize")

	if cityIdStr == "" {
		log.Println("Invalid cityId: ", cityIdStr)
		ginCtx.String(http.StatusBadRequest, "invalid cityId")
		return
	}
	cityId, err := uuid.Parse(cityIdStr)
	if err != nil {
		log.Println("Fail to parse id: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}

	filter := model.Filter{
		LocalID: cityId,
		PageSettings: model.PageSettings{
			Limit:  pageSize,
			Offset: (pageNum - 1) * pageSize,
		},
	}
	//TODO CREATE BODY
	_, _, err = h.useCase.GetByFilter(ctx, &filter)

	if err != nil {
		log.Println("Err to get offer by id: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := &docs.GetOffersResponse{}

	ginCtx.JSON(http.StatusOK, resp)
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
func (h *offerHandler) UpdateOffer(ctx *gin.Context) {
	var request docs.UpdateOfferRequest

	if err := ctx.BindJSON(&request); err != nil {
		log.Println("Invalid body")
		ctx.String(http.StatusBadRequest, "invalid body")
		return
	}

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)

	if idStr == "" || err != nil {
		log.Println("invalid offer id", idStr)
		ctx.String(http.StatusBadRequest, "invalid offer id")
		return
	}

	_ = id

	ctx.Status(http.StatusOK)
}
