package server

import "net/http"

// MapRoutes - maps the routes to the handlers
func (h *Handler) MapRoutes() {
	h.Router.HandleFunc("/ready", h.ReadyCheck).Methods(http.MethodGet)
	apiRouter := h.Router.Methods(http.MethodPost, http.MethodGet).PathPrefix("/api/v1/rates").Subrouter()
	apiRouter.HandleFunc("/latest", h.GetLatestRate).Queries("quote_currency", "{quote_currency}").Methods(http.MethodGet)
	apiRouter.HandleFunc("/timeseries", h.GetTimeseriesData).Queries("date", "{date}").Methods(http.MethodGet)
	apiRouter.HandleFunc("/range", h.GetRatesInRange).
		Queries(
			"quote_currency", "{quote_currency}",
			"from", "{from}",
			"to", "{to}",
		).
		Methods(http.MethodGet)

	apiRouter.HandleFunc("/{currency}", h.StoreRate).Methods(http.MethodPost)

	apiRouter.Use(JSONMiddleware)
}
