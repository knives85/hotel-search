package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer() *Server { return NewServer("/hotel-search") }

func TestHotelsIndex_NoParams_StillNotImplemented(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/hotel-search/hotels", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	// Parse + buildQuery succeed on empty input → handler falls through to 501.
	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotImplemented)
	}
}

func TestHotelsIndex_ValidParams_StillNotImplemented(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet,
		"/hotel-search/hotels?hotelName=Hilton&starRatings=FOUR&countryCodes=IT&page=2",
		nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusNotImplemented, rec.Body.String())
	}
}

func TestHotelsIndex_InvalidUniqueID_400(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/hotel-search/hotels?uniqueId=abc", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "uniqueId") {
		t.Errorf("body = %q, want to mention uniqueId", rec.Body.String())
	}
}

func TestHotelsIndex_InvalidSellStatus_400(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/hotel-search/hotels?sellStatus=yes", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHotelsIndex_InvalidDate_400(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet,
		"/hotel-search/hotels?creationDateFrom=15-01-2024", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if !strings.Contains(strings.ToLower(rec.Body.String()), "date") {
		t.Errorf("body = %q, want to mention date format", rec.Body.String())
	}
}
