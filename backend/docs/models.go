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

type GetUserAppLimitInfoResponse struct {
	Limit          uint `json:"limit"`
	ActiveAppCount uint `json:"active_app_count"`
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
	ID                string    `json:"id"`
	Task              string    `json:"task"`
	RoomID            string    `json:"room_id"`
	RoomName          string    `json:"room_name"`
	HotelID           string    `json:"hotel_id"`
	HotelName         string    `json:"hotel_name"`
	CheckIn           time.Time `json:"check_in_at"`
	CheckOut          time.Time `json:"check_out_at"`
	ExpirationAt      time.Time `json:"expiration_at"`
	ParticipantsLimit uint      `json:"participants_limit"`
	ParticipantsCount uint      `json:"participants_count"`
}

type CreateOfferRequest struct {
	HotelId           string    `json:"hotel_id" binding:"required"`
	LocationId        string    `json:"location_id" binding:"required"`
	ExpirationAt      time.Time `json:"expiration_at" binding:"required"`
	Task              string    `json:"task" binding:"required"`
	ParticipantsLimit uint      `json:"participants_limit" binding:"required"`
	RoomID            string    `json:"room_id" binding:"required"`
	CheckIn           time.Time `json:"check_in" binding:"required"`
	CheckOut          time.Time `json:"check_out" binding:"required"`
}

type CreateOfferResponse struct {
	Id string `json:"id"`
}

type GetOffersResponse struct {
	Offers     []*OfferResponse `json:"offers"`
	PagesCount int              `json:"pages_count"`
}

type UpdateOfferRequest struct {
	Task         string    `json:"task"`
	RoomID       string    `json:"room_id"`
	HotelID      string    `json:"hotel_id"`
	CheckIn      time.Time `json:"check_in_at"`
	CheckOut     time.Time `json:"check_out_at"`
	ExpirationAT time.Time `json:"expiration_at"`
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
	Rating        int    `json:"rating"`
}

func UserModelToResponse(u *user.User) *UserResponse {
	return &UserResponse{
		Id:            u.ID.String(),
		OstrovokLogin: u.OstrovokLogin,
		Email:         u.Email,
		IsAdmin:       u.IsAdmin,
		Rating:        u.Rating,
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
	UserId       string                 `json:"user_id"`
	HotelName    string                 `json:"hotel_name"`
	LocationName string                 `json:"location_name"`
	RoomName     string                 `json:"room_name"`
	Task         string                 `json:"task"`
	CheckInAt    time.Time              `json:"check_in_at"`
	CheckOutAt   time.Time              `json:"check_out_at"`
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

type GetByApplicationIdResponse struct {
	Id string `json:"id"`
}

type AnalyticsResponse struct {
	CompletedOffers      uint64 `json:"completed_offers"`
	ApplicationsReceived uint64 `json:"applications_received"`
	AcceptedReports      uint64 `json:"accepted_reports"`
}
