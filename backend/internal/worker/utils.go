package worker

import (
	"errors"
	"math/rand"

	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
)

func SelectRandomApplication(applications []*application.Application) (*application.Application, error) {
	if len(applications) == 0 {
		return nil, errors.New("no applications found")
	}

	randomIndex := rand.Intn(len(applications))
	return applications[randomIndex], nil
}
