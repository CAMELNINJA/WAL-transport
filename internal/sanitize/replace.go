package sanitize

import (
	"github.com/CAMELNINGA/cdc-postgres.git/internal/models"
)

// ReplaceHandler is a handler that replaces certain values in the request with new ones
type ReplaceHandler struct {
	BaseHandler
	oldTable       string
	newTable       string
	replaseColumns map[string]string
}

// Handle checks if the request contains the old value and replaces it with the new one if it does
func (h *ReplaceHandler) Handle(request models.ActionData) models.ActionData {
	if h.oldTable != "" || h.oldTable == request.Table || h.oldTable == "*" {
		request.Table = h.newTable
		if h.replaseColumns != nil {
			for _, v := range request.NewColumns {
				_, prs := h.replaseColumns[v.Name]
				if prs {
					v.Value = h.replaseColumns[v.Name]
				}
			}
		}
	}
	return h.BaseHandler.Handle(request)
}