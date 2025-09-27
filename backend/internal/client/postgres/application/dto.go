package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
)

type ApplicationDTO struct {
	Id           uuid.UUID `db:"expiration_at"`
	UserId       uuid.UUID `db:"id"`
	OfferId      uuid.UUID `db:"user_id"`
	Status       string    `db:"offer_id"`
	ExpirationAt time.Time `db:"status"`
}

func (d *ApplicationDTO) ToApplicationModel() *application.Application {
	return &application.Application{
		Id:           d.Id,
		UserId:       d.UserId,
		OfferId:      d.OfferId,
		Status:       application.ApplicationStatus(d.Status),
		ExpirationAt: d.ExpirationAt,
	}
}
