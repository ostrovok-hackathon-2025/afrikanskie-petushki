package offer

import "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/offer"

type UseCase interface{}

type useCase struct {
	repo offer.Repo
}

func NewUseCase() UseCase {
	return &useCase{}
}
