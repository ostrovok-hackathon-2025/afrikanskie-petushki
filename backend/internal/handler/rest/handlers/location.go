package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/location"
)

type LocationHandler interface {
	CreateLocation(ctx *gin.Context)
	GetLocations(ctx *gin.Context)
}

type locationHandler struct {
	useCase location.UseCase
}

func NewLocationHandler(useCase location.UseCase) LocationHandler {
	return &locationHandler{
		useCase: useCase,
	}
}

// CreateLocation
// Add godoc
// @Summary Create location
// @Description Creates location with given info
// @Tags Location
// @Accept json
// @Param input body docs.CreateLocationRequest true "Data for creating location"
// @Produce json
// @Security BearerAuth
// @Success 201 {object} docs.CreateLocationResponse "Created location data"
// @Failure 400 {string} string "Invalid data for creating location"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 500 "Internal server error"
// @Router /location/ [post]
func (h *locationHandler) CreateLocation(ginCtx *gin.Context) {
	var request docs.CreateLocationRequest
	ctx := context.Background()
	if err := ginCtx.BindJSON(&request); err != nil {
		log.Println("Invalid body")
		ginCtx.String(http.StatusBadRequest, "invalid body")
		return
	}

	id, err := h.useCase.Create(ctx, request.Name)
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

// GetLocations
// Add godoc
// @Summary Get locations
// @Description GetLocations all locations
// @Tags Location
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.GetLocationsResponse "Page of locations"
// @Failure 400 {string} string "Invalid data for getting locations"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 404 "Page with given number not found"
// @Failure 500 "Internal server error"
// @Router /location/ [get]
func (h *locationHandler) GetLocations(ginCtx *gin.Context) {
	ctx := context.Background()
	ucLocations, err := h.useCase.GetAll(ctx)
	if err != nil {
		log.Println("Err to get locations: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}
	apiLocations := make([]*docs.LocationResponse, len(ucLocations))
	for i, ucLocation := range ucLocations {
		apiLocations[i] = &docs.LocationResponse{
			Id:   ucLocation.ID.String(),
			Name: ucLocation.Name,
		}
	}

	resp := &docs.GetLocationsResponse{
		Locations: apiLocations,
	}

	ginCtx.JSON(http.StatusOK, resp)
}
