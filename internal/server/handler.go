package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	database "github.com/Shambou/golang-challenge/internal/database/postgres"
	"github.com/Shambou/golang-challenge/internal/seeds"
	"github.com/gorilla/mux"
)

type Handler struct {
	Router *mux.Router
	Server *http.Server
	DB     *database.Database
}

type JsonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// New - creates a new HTTP handler
func New() *Handler {
	h := &Handler{
		DB: database.NewDatabase(),
	}

	h.Router = mux.NewRouter()
	h.MapRoutes()

	err := h.DB.MigrateDB()
	if err != nil && err.Error() != "no change" {
		log.Println("failed to setup database", err)
	}

	seeder := seeds.New(h.DB)
	seeder.Execute()

	h.Server = &http.Server{
		Addr:    "0.0.0.0:" + os.Getenv("PORT"),
		Handler: h.Router,
		// Good practice to set timeouts to avoid slow loris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	return h
}

// Serve - gracefully serves our newly set up handler function
func (h *Handler) Serve() error {
	log.Printf("Serving app on :%s port", os.Getenv("PORT"))
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := h.Server.Shutdown(ctx)
	if err != nil {
		return err
	}

	log.Println("Shutting down gracefully")
	return nil
}

// ReadyCheck - Check if are connected to the database
func (h *Handler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("I am Ready!"); err != nil {
		panic(err)
	}
}

// jsonResponse - renders json response
func jsonResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(JsonResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}); err != nil {
		panic(err)
	}
}
