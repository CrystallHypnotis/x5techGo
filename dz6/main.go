package main

import (
	"context"
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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(status int) {
	if sr.status != 0 {
		return
	}
	sr.status = status
	sr.ResponseWriter.WriteHeader(status)
}

func (sr *statusRecorder) Write(body []byte) (int, error) {
	if sr.status == 0 {
		sr.WriteHeader(http.StatusOK)
	}
	return sr.ResponseWriter.Write(body)
}

func newRequestID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

type requestIDKey struct{}

func requestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

func requestLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID, err := newRequestID()
			if err != nil {
				logger.Error("fail when try to gen requestID", slog.String("error", err.Error()))
			}
			w.Header().Set("X-Request-ID", requestID)

			recorder := &statusRecorder{ResponseWriter: w}
			next.ServeHTTP(recorder, r)
			status := recorder.status
			if status == 0 {
				status = http.StatusOK
			}

			logger.Info("request received",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", status),
				slog.Int64("duration_ms", time.Since(start).Milliseconds()),
				slog.String("request_id", requestID))
		})
	}
}

func recoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("panic recovered",
						slog.Any("panic", recovered),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
						slog.String("request_id", requestIDFromContext(r.Context())))
					writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

type Handler struct {
	store *Store
}

func NewRouter(dataStore *Store, logger *slog.Logger) http.Handler {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
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
	if !isJSONContentType(r) {
		writeJSON(w, http.StatusUnsupportedMediaType, map[string]string{"error": "Content-Type must be application/json"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	var input struct {
		Customer string `json:"customer"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeJSON(w, decodeErrorStatus(err), map[string]string{"error": "invalid JSON"})
		return
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeJSON(w, decodeErrorStatus(err), map[string]string{"error": "request must contain a single JSON value"})
		return
	}

	customer := strings.TrimSpace(input.Customer)
	if customer == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "customer is required"})
		return
	}
	created := h.store.Create(customer)
	w.Header().Set("Location", fmt.Sprintf("/orders/%d", created.ID))
	writeJSON(w, http.StatusCreated, created)
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

func isJSONContentType(r *http.Request) bool {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	return err == nil && strings.EqualFold(contentType, "application/json")
}

func decodeErrorStatus(err error) int {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return http.StatusRequestEntityTooLarge
	}
	return http.StatusBadRequest
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
