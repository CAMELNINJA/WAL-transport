package sanitize

import (
	"github.com/CAMELNINGA/WAL-transport.git/internal/models"
)

// ReplaceHandler is a handler that replaces certain values in the request with new ones
type ReplaceHandler struct {
	BaseHandler
	oldTable       string
	newTable       string
	Schema         map[string]string
	replaseColumns map[string]string
}

// Handle checks if the request contains the old value and replaces it with the new one if it does
func (h *ReplaceHandler) Handle(request *models.ActionData) *models.ActionData {
	if h.oldTable == request.Table || h.oldTable == "*" {
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
	if h.BaseHandler.nextHandler != nil {
		return h.BaseHandler.Handle(request)
	}
	return request
}

type ReplaceOpts func(*ReplaceHandler)

func NewReplaceHandler(opts ...ReplaceOpts) Handler {
	handler := &ReplaceHandler{}
	for _, opt := range opts {
		opt(handler)
	}
	return handler
}

func WithReplaceTable(oldTable string, newTable string) ReplaceOpts {
	return func(handler *ReplaceHandler) {
		handler.oldTable = oldTable
		handler.newTable = newTable
	}
}

func WithReplaceColumns(replaseColumns map[string]string) ReplaceOpts {
	return func(handler *ReplaceHandler) {
		handler.replaseColumns = replaseColumns
	}
}

func WithReplaceSchema(schema map[string]string) ReplaceOpts {
	return func(handler *ReplaceHandler) {
		handler.Schema = schema
	}
}
