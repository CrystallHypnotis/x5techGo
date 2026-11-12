package handler

import (
	"log/slog"
	"net/http"
	"time"

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

func requestLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: capture status code and duration, generate request_id,
			// and write method, path, status, duration and request_id to slog.
			startMoment := time.Now()
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
				slog.Duration("duration", time.Since(startMoment)),
				slog.String("request_id", requestID))
		})
	}
}

func recoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: recover from panic, log the incident and return a safe JSON 500 response.
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("panic recovered",
						slog.Any("panic", recovered),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
					)

					writeJSON(w, http.StatusInternalServerError, errorResponse{
						Error: "internal error",
						Code:  "INTERNAL_SERVER_ERROR",
					})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
