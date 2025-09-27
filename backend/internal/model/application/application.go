package application

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationStatus string

const (
	APPLICATION_CREATED  = ApplicationStatus("__application_created")
	APPLICATION_ACCEPTED = ApplicationStatus("__application_accepted")
	APPLICATION_DECLINED = ApplicationStatus("__application_declined")
)

type Application struct {
	Id           uuid.UUID
	UserId       uuid.UUID
	OfferId      uuid.UUID
	Status       ApplicationStatus
	ExpirationAt time.Time
}

func NewApplication(userId, offerId uuid.UUID) *Application {
	return &Application{
		Id:      uuid.New(),
		UserId:  userId,
		OfferId: offerId,
		Status:  APPLICATION_CREATED,
	}
}
