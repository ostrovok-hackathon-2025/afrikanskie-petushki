package dto

type TokenRequest struct {
	Auth string `header:"Authorization" binding:"required"`
}
