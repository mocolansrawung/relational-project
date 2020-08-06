package handlers

import (
	"github.com/evermos/boilerplate-go/src/services"
)

type Handler struct {
	Service *services.Service `inject:"service"`
}

func (h *Handler) TestHandler() {
	h.Service.TestService()
}
