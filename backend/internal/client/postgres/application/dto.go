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
