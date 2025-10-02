package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/hotel"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/hotel"
)

type HotelHandler interface {
	CreateHotel(ctx *gin.Context)
	GetHotels(ctx *gin.Context)
}

type hotelHandler struct {
	useCase hotel.UseCase
}

func NewHotelHandler(useCase hotel.UseCase) HotelHandler {
	return &hotelHandler{
		useCase: useCase,
	}
}

// CreateHotel
// Add godoc
// @Summary Create hotel
// @Description Creates offer with given info
// @Tags Hotel
// @Accept json
// @Param input body docs.CreateHotelRequest true "Data for creating hotel"
// @Produce json
// @Security BearerAuth
// @Success 201 {object} docs.CreateHotelResponse "Created hotel data"
// @Failure 400 {string} string "Invalid data for creating hotel"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 500 "Internal server error"
// @Router /hotel/ [post]
func (h *hotelHandler) CreateHotel(ginCtx *gin.Context) {
	var request docs.CreateHotelRequest
	ctx := context.Background()
	if err := ginCtx.BindJSON(&request); err != nil {
		log.Println("Invalid body")
		ginCtx.String(http.StatusBadRequest, "invalid body")
		return
	}

	locationIDStr := request.LocationID
	locationID, err := uuid.Parse(locationIDStr)
	if locationIDStr == "" || err != nil {
		log.Println("Invalid location_id: ", locationIDStr)
		ginCtx.String(http.StatusBadRequest, "invalid location_id")
		return
	}

	create := model.Create{Name: request.Name, LocationID: locationID}
	id, err := h.useCase.Create(ctx, create)
	if err != nil {
		log.Println("Err to create location: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}
	resp := &docs.CreateOfferResponse{
		Id: id.String(),
	}

	ginCtx.JSON(http.StatusCreated, resp)
}

// GetHotels
// Add godoc
// @Summary Get hotels
// @Description GetHotels all hotels
// @Tags Hotel
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.GetHotelsResponse "Page of hotels"
// @Failure 400 {string} string "Invalid data for getting hotels"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 404 "Page with given number not found"
// @Failure 500 "Internal server error"
// @Router /hotel/ [get]
func (h *hotelHandler) GetHotels(ginCtx *gin.Context) {
	ctx := context.Background()
	ucHotels, err := h.useCase.GetAll(ctx)
	if err != nil {
		log.Println("Err to get hotels: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}
	apiHotels := make([]*docs.HotelResponse, len(ucHotels))
	for i, ucHotel := range ucHotels {
		apiHotels[i] = &docs.HotelResponse{
			Id:           ucHotel.ID.String(),
			Name:         ucHotel.Name,
			LocationId:   ucHotel.LocationID.String(),
			LocationName: ucHotel.LocationName,
		}
	}
	ginCtx.JSON(http.StatusOK, docs.GetHotelsResponse{Hotels: apiHotels})
}
