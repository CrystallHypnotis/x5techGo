package main

import (
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
		t.Fatalf("status = %d, expected 200", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type = %q", rec.Header().Get("Content-Type"))
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"status":"ok"}` {
		t.Fatalf(`body = %q, expected {"status":"ok"}`, got)
	}
}

// Добавьте проверки успешного создания и ошибочного входа.
func TestCreateOrder(t *testing.T) {
	router := NewRouter(NewStore(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"customer":" Ded Moroz    "}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, expected 201", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type = %q", rec.Header().Get("Content-Type"))
	}
	if rec.Header().Get("Location") != "/orders/1" {
		t.Fatalf("Location = %q, expected /orders/1", rec.Header().Get("Location"))
	}
	var order Order
	if err := json.Unmarshal(rec.Body.Bytes(), &order); err != nil {
		t.Fatal(err)
	}
	if order.ID != 1 || order.Customer != "Ded Moroz" || order.Status != "new" {
		t.Fatalf("created order = %+v", order)
	}
}

func TestCreateOrderInvalidInput(t *testing.T) {
	store := NewStore()
	router := NewRouter(store, slog.New(slog.NewTextHandler(io.Discard, nil)))
	tests := []string{"", `{"customer":`, `{"customer":9999}`, `{"customer":""}`}
	for _, test := range tests {
		req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(test))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("test = %q: status = %d, expected 400", test, rec.Code)
		}
		if rec.Header().Get("Content-Type") != "application/json" {
			t.Errorf("test = %q: Content-Type = %q", test, rec.Header().Get("Content-Type"))
		}
		var response map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil || response["error"] == "" {
			t.Errorf("test = %q: invalid error response %s", test, rec.Body.String())
		}
	}
	if _, ok := store.ByID(1); ok {
		t.Fatal("invalid input created order")
	}
}
