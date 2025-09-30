package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
)

type LimitsDTO struct {
	ParticipantsLimit uint `db:"participants_limit"`
	ParticipantsCount uint `db:"participants_count"`
}

type UserAppLimitInfoDTO struct {
	Limit          uint `db:"app_limit"`
	ActiveAppCount uint `db:"active_app_count"`
}

type ApplicationDTO struct {
	Id           uuid.UUID `db:"id"`
	UserId       uuid.UUID `db:"user_id"`
	OfferId      uuid.UUID `db:"offer_id"`
	Status       string    `db:"status"`
	ExpirationAt time.Time `db:"expiration_at"`
	HotelName    string    `db:"name"`
}

func (d *ApplicationDTO) ToApplicationModel() *application.Application {
	return &application.Application{
		Id:           d.Id,
		UserId:       d.UserId,
		OfferId:      d.OfferId,
		Status:       application.ApplicationStatus(d.Status),
		ExpirationAt: d.ExpirationAt,
		HotelName:    d.HotelName,
	}
}

func (d *UserAppLimitInfoDTO) ToModel() *application.UserAppLimitInfo {
	return &application.UserAppLimitInfo{
		Limit:          d.Limit,
		ActiveAppCount: d.ActiveAppCount,
	}
}
