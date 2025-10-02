package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/pkg"
)

type ApplicationStatus string

const (
	APPLICATION_CREATED  = ApplicationStatus("__app_created")
	APPLICATION_ACCEPTED = ApplicationStatus("__app_accepted")
	APPLICATION_DECLINED = ApplicationStatus("__app_declined")
)

type Application struct {
	Id           uuid.UUID
	UserId       uuid.UUID
	OfferId      uuid.UUID
	Status       ApplicationStatus
	ExpirationAt time.Time
	HotelName    string
}

type UserAppLimitInfo struct {
	Limit          uint
	ActiveAppCount uint
}

func NewApplication(userId, offerId uuid.UUID) *Application {
	return &Application{
		Id:      uuid.New(),
		UserId:  userId,
		OfferId: offerId,
		Status:  APPLICATION_CREATED,
	}
}

type Filter struct {
	HotelID       pkg.Opt[uuid.UUID]
	LocationID    pkg.Opt[uuid.UUID]
	RoomID        pkg.Opt[uuid.UUID]
	Status        pkg.Opt[ApplicationStatus]
	Limit, Offset uint64
}
