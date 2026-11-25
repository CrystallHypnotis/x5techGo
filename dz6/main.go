package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type CreateOrderInput struct {
	Customer string `json:"customer"`
}

type Handler struct {
	store *Store
}

func NewRouter(dataStore *Store, logger *slog.Logger) http.Handler {
	h := &Handler{store: dataStore}
	router := chi.NewRouter()
	router.Use(requestLoggingMiddleware(logger))
	router.Use(recoveryMiddleware(logger))

	router.Get("/health", h.health)
	router.Post("/orders", h.createOrder)
	router.Get("/orders/{id}", h.getOrder)
	return router
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || !strings.EqualFold(contentType, "application/json") {
		writeJSON(w, http.StatusUnsupportedMediaType,
			map[string]string{"error": "incorrect Content-Type"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1024)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input CreateOrderInput
	if err := decoder.Decode(&input); err != nil {
		if _, ok := err.(*http.MaxBytesError); ok {
			writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "request too large"})
		} else {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		}
		return
	}

	var extra CreateOrderInput
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if _, ok := err.(*http.MaxBytesError); ok {
			writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "request too large"})
		} else {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		}
		return
	}

	input.Customer = strings.TrimSpace(input.Customer)
	if input.Customer == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty name"})
		return
	}
	newOrder := h.store.Create(input.Customer)
	w.Header().Set("Location", fmt.Sprintf("/orders/%d", newOrder.ID))
	writeJSON(w, http.StatusCreated, newOrder)
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid order id"})
		return
	}
	result, ok := h.store.ByID(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "order not found"})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		status = http.StatusInternalServerError
		data = []byte(`{"error":"internal server error"}`)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(append(data, '\n'))
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	router := NewRouter(NewStore(), logger)
	logger.Info("starting HTTP server", slog.String("address", ":8080"))

	server := &http.Server{Addr: ":8080", Handler: router}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("HTTP server stopped", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
