package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/docs"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/usecase/room"
)

type RoomHandler interface {
	// CreateRoom
	// Add godoc
	// @Summary Create offer
	// @Description Creates offer with given info
	// @Tags Offer
	// @Accept json
	// @Param input body docs.CreateRoomRequest true "Data for creating room"
	// @Produce json
	// @Security BearerAuth
	// @Success 201 {object} docs.CreateRoomResponse "Created room data"
	// @Failure 400 {string} string "Invalid data for creating room"
	// @Failure 401 "Unauthorized"
	// @Failure 403 "Only available for admin"
	// @Failure 500 "Internal server error"
	// @Router /room/ [post]
	CreateRoom(ctx *gin.Context)
	GetRooms(ctx *gin.Context)
}

type roomHandler struct {
	useCase room.UseCase
}

func NewRoomHandler(useCase room.UseCase) RoomHandler {
	return &roomHandler{
		useCase: useCase,
	}
}

// CreateRoom
// Add godoc
// @Summary Create offer
// @Description Creates offer with given info
// @Tags Room
// @Accept json
// @Param input body docs.CreateRoomRequest true "Data for creating room"
// @Produce json
// @Security BearerAuth
// @Success 201 {object} docs.CreateRoomResponse "Created room data"
// @Failure 400 {string} string "Invalid data for creating room"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 500 "Internal server error"
// @Router /room/ [post]
func (h *roomHandler) CreateRoom(ginCtx *gin.Context) {
	var request docs.CreateRoomRequest
	ctx := context.Background()
	if err := ginCtx.BindJSON(&request); err != nil {
		log.Println("Invalid body")
		ginCtx.String(http.StatusBadRequest, "invalid body")
		return
	}

	id, err := h.useCase.Create(ctx, request.Name)
	if err != nil {
		log.Println("Err to create room: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}
	resp := &docs.CreateOfferResponse{
		Id: id.String(),
	}

	ginCtx.JSON(http.StatusCreated, resp)
}

// GetRooms
// Add godoc
// @Summary Get rooms
// @Description GetRooms all rooms
// @Tags Room
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.GetRoomsResponse "Page of rooms"
// @Failure 400 {string} string "Invalid data for getting rooms"
// @Failure 401 "Unauthorized"
// @Failure 403 "Only available for admin"
// @Failure 404 "Page with given number not found"
// @Failure 500 "Internal server error"
// @Router /room/ [get]
func (h *roomHandler) GetRooms(ginCtx *gin.Context) {
	ctx := context.Background()
	ucRooms, err := h.useCase.GetAll(ctx)
	if err != nil {
		log.Println("Err to get rooms: ", err.Error())
		ginCtx.String(http.StatusBadRequest, err.Error())
		return
	}
	apiRooms := make([]*docs.RoomResponse, len(ucRooms))
	for i, ucRoom := range ucRooms {
		apiRooms[i] = &docs.RoomResponse{
			Id:   ucRoom.ID.String(),
			Name: ucRoom.Name,
		}
	}
	ginCtx.JSON(http.StatusOK, docs.GetRoomsResponse{Rooms: apiRooms})
}
