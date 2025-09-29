package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/offer"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/pkg"

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
func (h *offerHandler) CreateOffer(ctx *gin.Context) {
	var request docs.CreateOfferRequest

	if err := ctx.BindJSON(&request); err != nil {
		log.Println("Invalid body")
		ctx.String(http.StatusBadRequest, "invalid body")
		return
	}

	create := model.Create{}

	id, err := h.useCase.Create(ctx, create)

	if err != nil {
		log.Println("Err to create offer: ", err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := &docs.CreateOfferResponse{
		Id: id.String(),
	}

	ctx.JSON(http.StatusCreated, resp)
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
func (h *offerHandler) GetOffers(ctx *gin.Context) {
	pageNumStr := ctx.Query("pageNum")
	pageNum, err := strconv.ParseUint(pageNumStr, 10, 0)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageNum: ", pageNumStr)
		ctx.String(http.StatusBadRequest, "invalid pageNum")
		return
	}

	pageSizeStr := ctx.Query("pageSize")
	pageSize, err := strconv.ParseUint(pageSizeStr, 10, 0)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageSize: ", pageSizeStr)
		ctx.String(http.StatusBadRequest, "invalid pageSize")
		return
	}
	pageSettings := model.PageSettings{
		Limit:  pageSize,
		Offset: pageNum * pageSize,
	}
	//TODO CREATE BODY
	ucOffers, pagesCount, err := h.useCase.GetForPage(ctx, pageSettings)

	if err != nil {
		log.Println("Err to get offers for page: ", err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := &docs.GetOffersResponse{
		Offers:     convertUcOffersToApi(ucOffers),
		PagesCount: pagesCount,
	}

	ctx.JSON(http.StatusOK, resp)
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
func (h *offerHandler) GetOfferById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)

	if idStr == "" || err != nil {
		log.Println("invalid offer id", idStr)
		ctx.String(http.StatusBadRequest, "invalid offer id")
		return
	}

	ucOffer, err := h.useCase.GetByID(ctx, id)
	if err != nil {
		log.Println("Err to get offer by id: ", err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := convertUcOfferToApi(ucOffer)

	ctx.JSON(http.StatusOK, resp)
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
func (h *offerHandler) FindOffers(ctx *gin.Context) {
	pageNumStr := ctx.Query("pageNum")
	pageNum, err := strconv.ParseUint(pageNumStr, 10, 0)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageNum: ", pageNumStr)
		ctx.String(http.StatusBadRequest, "invalid pageNum")
		return
	}

	pageSizeStr := ctx.Query("pageSize")
	pageSize, err := strconv.ParseUint(pageSizeStr, 10, 0)

	if pageNumStr == "" || err != nil {
		log.Println("Invalid pageSize: ", pageSizeStr)
		ctx.String(http.StatusBadRequest, "invalid pageSize")
		return
	}

	cityIdStr := ctx.Query("pageSize")

	if cityIdStr == "" {
		log.Println("Invalid cityId: ", cityIdStr)
		ctx.String(http.StatusBadRequest, "invalid cityId")
		return
	}
	cityId, err := uuid.Parse(cityIdStr)
	if err != nil {
		log.Println("Fail to parse id: ", err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	filter := model.Filter{
		LocationID: pkg.NewWithValue(cityId),
		Limit:      pkg.NewWithValue(pageSize),
		Offset:     pkg.NewWithValue(pageNum * pageSize),
	}
	ucOffers, pagesCount, err := h.useCase.GetByFilter(ctx, filter)

	if err != nil {
		log.Println("Err to find offers by filter: ", err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := &docs.GetOffersResponse{
		Offers:     convertUcOffersToApi(ucOffers),
		PagesCount: pagesCount,
	}

	ctx.JSON(http.StatusOK, resp)
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
	var (
		count int
		edit  model.Edit
	)
	edit.OfferID = id
	if request.Task != "" {
		edit.Task = pkg.NewWithValue(request.Task)
		count++
	}
	if request.RoomID != "" {
		roomID, err := uuid.Parse(request.RoomID)
		if err != nil {
			log.Println("invalid room id", idStr)
			ctx.String(http.StatusBadRequest, "invalid room id")
			return
		}
		edit.RoomID = pkg.NewWithValue(roomID)
		count++
	}
	if request.HotelID != "" {
		hotelID, err := uuid.Parse(request.HotelID)
		if err != nil {
			log.Println("invalid hotel id", idStr)
			ctx.String(http.StatusBadRequest, "invalid hotel id")
			return
		}
		edit.RoomID = pkg.NewWithValue(hotelID)
		count++
	}
	//TODO Validation
	edit.CheckIn = pkg.NewWithValue(request.CheckIn)
	edit.CheckOut = pkg.NewWithValue(request.CheckOut)
	edit.ExpirationAT = pkg.NewWithValue(request.ExpirationAT)

	err = h.useCase.Edit(ctx, edit)
	if err != nil {
		log.Println("Err to update offer by id: ", err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func convertUcOffersToApi(ucOffers []model.Offer) []*docs.OfferResponse {
	apiOffers := make([]*docs.OfferResponse, len(ucOffers))
	for i, ucOffer := range ucOffers {
		apiOffers[i] = convertUcOfferToApi(ucOffer)
	}
	return apiOffers
}

func convertUcOfferToApi(ucOffer model.Offer) *docs.OfferResponse {
	return &docs.OfferResponse{
		ID:           ucOffer.ID.String(),
		Task:         ucOffer.Task,
		RoomID:       ucOffer.RoomID.String(),
		RoomName:     ucOffer.RoomName,
		HotelID:      ucOffer.HotelID.String(),
		HotelName:    ucOffer.HotelName,
		CheckIn:      ucOffer.CheckIn,
		CheckOut:     ucOffer.CheckOut,
		ExpirationAt: ucOffer.ExpirationAt,
	}
}
