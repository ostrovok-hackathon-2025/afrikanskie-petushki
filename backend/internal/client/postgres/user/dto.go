package user

import (
	"github.com/google/uuid"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
)

type UserDTO struct {
	Id            uuid.UUID `db:"id"`
	OstrovokLogin string    `db:"ostrovok_login"`
	IsAdmin       bool      `db:"is_admin"`
	Rating        int       `db:"rating"`
}

func (d *UserDTO) ToUserModel() *model.User {
	return &model.User{
		ID:            d.Id,
		OstrovokLogin: d.OstrovokLogin,
		IsAdmin:       d.IsAdmin,
		Rating:        d.Rating,
	}
}
