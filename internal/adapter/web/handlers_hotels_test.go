package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/knives85/hotel-search/internal/domain"
)

type stubSearchPort struct {
	result domain.SearchResult
	err    error
}

func (s stubSearchPort) Search(ctx context.Context, q domain.HotelSearchQuery) (domain.SearchResult, error) {
	return s.result, s.err
}
func (s stubSearchPort) SidebarFilterCounts(context.Context, domain.HotelSearchQuery) (domain.SidebarFilterCounts, error) {
	return domain.SidebarFilterCounts{}, nil
}

func newTestServer() *Server { return NewServer("/hotel-search", Deps{Search: stubSearchPort{}}) }

func TestHotelsIndex_NoParams_ReturnsOk(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/hotel-search/hotels", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	// Parse + buildQuery succeed on empty input → handler falls through to 501.
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestHotelsIndex_ValidParams_ReturnsOk(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet,
		"/hotel-search/hotels?hotelName=Hilton&starRatings=FOUR&countryCodes=IT&page=2",
		nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusOK, rec.Body.String())
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
