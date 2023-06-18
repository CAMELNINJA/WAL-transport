package sanitize

import "github.com/CAMELNINGA/cdc-postgres.git/internal/models"

// Handler interface defines the method that each handler in the chain must implement
type Handler interface {
	SetNext(handler Handler)                              // Set the next handler in the chain
	Handle(request *models.ActionData) *models.ActionData // Handle the request
}

// BaseHandler is the base struct that implements the Handler interface
type BaseHandler struct {
	nextHandler Handler
}

// SetNext sets the next handler in the chain
func (h *BaseHandler) SetNext(handler Handler) {
	h.nextHandler = handler
}

// Handle passes the request to the next handler in the chain, if there is one
func (h *BaseHandler) Handle(request *models.ActionData) *models.ActionData {
	if h.nextHandler != nil {
		return h.nextHandler.Handle(request)
	}
	return request
}

func NewSanitizeHandler() Handler {
	return &BaseHandler{}
}
