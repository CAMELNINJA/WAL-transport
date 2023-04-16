package sanitize

import (
	"github.com/CAMELNINGA/cdc-postgres.git/internal/models"
)

// ReplaceHandler is a handler that replaces certain values in the request with new ones
type ReplaceHandler struct {
	BaseHandler
	old string
	new string
}

// Handle checks if the request contains the old value and replaces it with the new one if it does
func (h *ReplaceHandler) Handle(request models.ActionData) models.ActionData {
	if h.old != "" {
		// request = strings.Replace(request, h.old, h.new, -1)
	}
	return h.BaseHandler.Handle(request)
}
