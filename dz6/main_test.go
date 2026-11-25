package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// До реализации health этот тест должен быть красным.
func TestHealth(t *testing.T) {
	router := NewRouter(NewStore(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type=%q", rec.Header().Get("Content-Type"))
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Status != "ok" {
		t.Fatalf("status in body=%q", body.Status)
	}
}

// Добавьте проверки успешного создания и ошибочного входа.
func TestCreateAndGetOrder(t *testing.T) {
	router := NewRouter(NewStore(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"customer":" Анна "}`))
	request.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, request)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status=%d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("POST Content-Type=%q", rec.Header().Get("Content-Type"))
	}
	if rec.Header().Get("Location") != "/orders/1" {
		t.Fatalf("Location=%q, want /orders/1", rec.Header().Get("Location"))
	}
	var created Order
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID != 1 || created.Customer != "Анна" || created.Status != "new" {
		t.Fatalf("created order=%+v", created)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK || getRec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("GET status=%d Content-Type=%q", getRec.Code, getRec.Header().Get("Content-Type"))
	}
	var fetched Order
	if err := json.Unmarshal(getRec.Body.Bytes(), &fetched); err != nil {
		t.Fatal(err)
	}
	if fetched != created {
		t.Fatalf("fetched order=%+v, want %+v", fetched, created)
	}
}

func TestGetOrderErrors(t *testing.T) {
	router := NewRouter(NewStore(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, tc := range []struct {
		path string
		want int
	}{
		{path: "/orders/missing", want: http.StatusBadRequest},
		{path: "/orders/0", want: http.StatusBadRequest},
		{path: "/orders/-1", want: http.StatusBadRequest},
		{path: "/orders/999999999999999999999999999999999", want: http.StatusBadRequest},
		{path: "/orders/42", want: http.StatusNotFound},
	} {
		request := httptest.NewRequest(http.MethodGet, tc.path, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, request)
		if rec.Code != tc.want {
			t.Errorf("GET %s status=%d, want %d", tc.path, rec.Code, tc.want)
		}
		if rec.Header().Get("Content-Type") != "application/json" {
			t.Errorf("GET %s Content-Type=%q", tc.path, rec.Header().Get("Content-Type"))
		}
	}
}

func TestCreateOrderRejectsInvalidJSON(t *testing.T) {
	store := NewStore()
	router := NewRouter(NewStore(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, tc := range []struct {
		name string
		body string
	}{
		{name: "empty body", body: ""},
		{name: "malformed JSON", body: `{"customer":`},
		{name: "wrong customer type", body: `{"customer":123}`},
		{name: "empty customer", body: `{"customer":"  "}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status=%d, want 400; body=%s", rec.Code, rec.Body.String())
			}
			if rec.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("Content-Type=%q", rec.Header().Get("Content-Type"))
			}
			var body map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body["error"] == "" {
				t.Fatalf("invalid error response: %s", rec.Body.String())
			}
		})
	}
	if _, ok := store.ByID(1); ok {
		t.Fatal("invalid requests must not create an order")
	}
}

func TestRecoveryWrites500ToAccessLog(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	panicHandler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("test panic") })
	h := requestLoggingMiddleware(logger)(recoveryMiddleware(logger)(panicHandler))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/panic", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d, want 500", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type=%q", rec.Header().Get("Content-Type"))
	}
	if !strings.Contains(logs.String(), `"status":500`) {
		t.Fatalf("access log does not contain status 500: %s", logs.String())
	}
	requestID := rec.Header().Get("X-Request-ID")
	if requestID == "" || !strings.Contains(logs.String(), `"request_id":"`+requestID+`"`) {
		t.Fatalf("request ID is absent or differs from access log: id=%q, logs=%s", requestID, logs.String())
	}
}
