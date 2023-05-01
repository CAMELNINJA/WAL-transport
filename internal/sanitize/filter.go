package sanitize

import (
	"github.com/CAMELNINGA/cdc-postgres.git/internal/models"
)

// FilterHandler is a handler that filters out certain data from the request
type FilterHandler struct {
	BaseHandler
	Table         string
	Schema        map[string]string
	filterColumns map[string]string
}

// Handle checks if the request contains the filter and removes it if it does
// Todo filter tables and columns
func (h *FilterHandler) Handle(request *models.ActionData) *models.ActionData {

	if h.Table == request.Table || h.Table == "*" {
		if h.Schema != nil {
			_, prs := h.Schema[request.Schema]
			if prs {
				return nil
			}

		}
		if h.filterColumns != nil {

			for i, v := range request.NewColumns {
				_, prs := h.filterColumns[v.Name]
				if prs {
					request.NewColumns = append(request.NewColumns[:i], request.NewColumns[i+1:]...)
				}
			}

			for i, v := range request.OldColumns {
				_, prs := h.filterColumns[v.Name]
				if prs {
					request.OldColumns = append(request.OldColumns[:i], request.OldColumns[i+1:]...)
				}
			}

		} else if h.Schema != nil {
			return request
		} else {
			return nil
		}
	}
	if len(request.NewColumns) == 0 && len(request.OldColumns) == 0 {
		return nil
	}

	if h.BaseHandler.nextHandler != nil {
		return h.BaseHandler.Handle(request)
	}
	return request
}
