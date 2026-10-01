package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHealth(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/healthz", nil), rec)

	if err := Health(c); err != nil {
		t.Fatalf("Health returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestReady(t *testing.T) {
	e := echo.New()

	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/readyz", nil), rec)
	if err := Ready(func(context.Context) error { return nil })(c); err != nil {
		t.Fatalf("Ready returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 when dependency is healthy, got %d", rec.Code)
	}

	recDown := httptest.NewRecorder()
	cDown := e.NewContext(httptest.NewRequest(http.MethodGet, "/readyz", nil), recDown)
	if err := Ready(func(context.Context) error { return errors.New("db down") })(cDown); err != nil {
		t.Fatalf("Ready returned error: %v", err)
	}
	if recDown.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 when dependency is down, got %d", recDown.Code)
	}
}
