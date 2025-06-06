package handlers

import (
	"github.com/freekobie/hazel/services"
	"github.com/go-playground/validator/v10"
)

// TODO: remove global variable
var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type Handler struct {
	us  *services.UserService
	wss *services.WorkspaceService
}

func NewHandler(us *services.UserService, wks *services.WorkspaceService) *Handler {
	return &Handler{
		us:  us,
		wss: wks,
	}
}
