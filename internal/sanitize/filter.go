package sanitize

import "github.com/CAMELNINGA/cdc-postgres.git/internal/models"

// FilterHandler is a handler that filters out certain data from the request
type FilterHandler struct {
	BaseHandler
	filter string
}

// Handle checks if the request contains the filter and removes it if it does
func (h *FilterHandler) Handle(request models.ActionData) models.ActionData {
	// if h.filter != "" {
	// 	request = strings.Replace(request, h.filter, "", -1)
	// }
	return h.BaseHandler.Handle(request)
}
