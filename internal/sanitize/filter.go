package sanitize

import "github.com/CAMELNINGA/cdc-postgres.git/internal/models"

// FilterHandler is a handler that filters out certain data from the request
type FilterHandler struct {
	BaseHandler
	Table         string
	filterColumns FilterColumns
}

type FilterColumns []*FilterColumn
type FilterColumn struct {
	Column string
}

// Handle checks if the request contains the filter and removes it if it does
// Todo filter tables and columns
func (h *FilterHandler) Handle(request models.ActionData) models.ActionData {

	if h.Table != "" || h.Table == request.Table {
		if h.filterColumns != nil {
			for _, v := range h.filterColumns {
				for i, v2 := range request.NewColumns {
					if v.Column == v2.Name {
						request.NewColumns = append(request.NewColumns[:i], request.NewColumns[i+1:]...)
					}
				}
			}
		} else {
			request.NewColumns = nil
		}
	}
	return h.BaseHandler.Handle(request)
}
