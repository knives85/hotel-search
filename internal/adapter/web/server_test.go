package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	srv := NewServer("/hotel-search")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("healthz: got status %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != "ok" {
		t.Fatalf("healthz: got body %q, want %q", got, "ok")
	}
}

func TestHotelsRouteRegistered(t *testing.T) {
	srv := NewServer("/hotel-search")

	req := httptest.NewRequest(http.MethodGet, "/hotel-search/hotels", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	// The route exists but is not implemented yet: expect 501, not 404.
	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("/hotels: got status %d, want %d", rec.Code, http.StatusNotImplemented)
	}
}
