package server

import "net/http"

// MapRoutes - maps the routes to the handlers
func (h *Handler) MapRoutes() {
	h.Router.HandleFunc("/ready", h.ReadyCheck).Methods(http.MethodGet)
	
}
