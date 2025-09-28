package docs

import (
	"mime/multipart"
	"time"

	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
)

type ApplicationResponse struct {
	Id           string    `json:"id"`
	UserId       string    `json:"user_id"`
	OfferId      string    `json:"offer_id"`
	Status       string    `json:"status"`
	ExpirationAt time.Time `json:"expiration_at"`
	HotelName    string    `json:"hotel_name"`
}

func ApplicationModelToResponse(model *application.Application) *ApplicationResponse {
	return &ApplicationResponse{
		Id:           model.Id.String(),
		UserId:       model.UserId.String(),
		OfferId:      model.OfferId.String(),
		Status:       string(model.Status),
		ExpirationAt: model.ExpirationAt,
		HotelName:    model.HotelName,
	}
}

type CreateApplicationRequest struct {
	OfferId string `json:"offer_id" binding:"required"`
}

type CreateApplicationResponse struct {
	ApplicationId string `json:"application_id"`
}

type GetApplicationsResponse struct {
	Applications []*ApplicationResponse `json:"applications"`
	PagesCount   int                    `json:"pages_count"`
}

type OfferResponse struct {
	Id           string    `json:"id"`
	HotelId      string    `json:"hotel_id"`
	LocationId   string    `json:"location_id"`
	ExpirationAt time.Time `json:"expiration_at"`
	CheckDate    time.Time `json:"check_date"`
	Task         string    `json:"task"`
	Used         bool      `json:"used"`
}

type CreateOfferRequest struct {
	HotelId      string    `json:"hotel_id" binding:"required"`
	LocationId   string    `json:"location_id" binding:"required"`
	ExpirationAt time.Time `json:"expiration_at" binding:"required"`
	Task         string    `json:"task" binding:"required"`
}

type CreateOfferResponse struct {
	Id string `json:"id"`
}

type GetOffersResponse struct {
	Offers     []*OfferResponse `json:"offers"`
	PagesCount int              `json:"pages_count"`
}

type UpdateOfferRequest struct {
	ExpirationAt time.Time `json:"expiration_at"`
	Task         string    `json:"task"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessTTL    int    `json:"access_ttl"`
	RefreshTTL   int    `json:"refresh_ttl"`
}

type UserResponse struct {
	Id            string `json:"id"`
	OstrovokLogin string `json:"ostrovok_login"`
	Email         string `json:"email"`
	IsAdmin       bool   `json:"is_admin"`
}

func UserModelToResponse(u *user.User) *UserResponse {
	return &UserResponse{
		Id:            u.ID.String(),
		OstrovokLogin: u.OstrovokLogin,
		Email:         u.Email,
		IsAdmin:       u.IsAdmin,
	}
}

type LogInRequest struct {
	OstrovokLogin string `json:"ostrovok_login" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

type SignUpRequest struct {
	OstrovokLogin string `json:"ostrovok_login" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ReportImageResponse struct {
	Id   string `json:"id"`
	Link string `json:"link"`
}

type ReportResponse struct {
	Id           string                 `json:"id"`
	ExpirationAt string                 `json:"expiration_at"`
	Status       string                 `json:"status"`
	Text         string                 `json:"text"`
	Images       []*ReportImageResponse `json:"images"`
}

type GetReportsResponse struct {
	Reports    []*ReportResponse `json:"reports"`
	PagesCount int               `json:"pages_count"`
}

type UpdateReportRequest struct {
	Text   string                  `form:"text"`
	Images []*multipart.FileHeader `form:"image"`
}

type ConfirmReport struct {
	Status string `json:"status" binding:"required"`
}

type CreateHotelRequest struct {
	Name       string `json:"name" binding:"required"`
	LocationID string `json:"location_id" binding:"required"`
}

type CreateHotelResponse struct {
	HotelId string `json:"hotel_id"`
}

type CreateLocationRequest struct {
	Name string `json:"name" binding:"required"`
}

type HotelResponse struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	LocationId   string `json:"location_id"`
	LocationName string `json:"location_name"`
}

type GetHotelsResponse struct {
	Hotels []*HotelResponse `json:"hotels"`
}

type CreateLocationResponse struct {
	LocationId string `json:"location_id"`
}

type LocationResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetLocationsResponse struct {
	Locations []*LocationResponse `json:"locations"`
}

type CreateRoomRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateRoomResponse struct {
	RoomId string `json:"room_id"`
}

type RoomResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetRoomsResponse struct {
	Rooms []*RoomResponse `json:"rooms"`
}
