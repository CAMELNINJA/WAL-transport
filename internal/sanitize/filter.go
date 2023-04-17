package sanitize

import (
	"github.com/CAMELNINGA/cdc-postgres.git/internal/models"
)

// FilterHandler is a handler that filters out certain data from the request
type FilterHandler struct {
	BaseHandler
	Table         string
	filterColumns map[string]string
}

// Handle checks if the request contains the filter and removes it if it does
// Todo filter tables and columns
func (h *FilterHandler) Handle(request models.ActionData) models.ActionData {

	if h.Table != "" || h.Table == request.Table {
		if h.filterColumns != nil {

			for i, v := range request.NewColumns {
				_, prs := h.filterColumns[v.Name]
				if prs {
					request.NewColumns = append(request.NewColumns[:i], request.NewColumns[i+1:]...)
				}
			}

		} else {
			request.NewColumns = nil
		}
	}
	return h.BaseHandler.Handle(request)
}