package dto

type AuthenticateThroughTelergamRequest struct {
	ID        string `query:"id" validate:"required"`
	FirstName string `query:"first_name" validate:"required"`
	LastName  string `query:"last_name" validate:"required"`
	PhotoURL  string `query:"photo_url"`
}
