package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)


// -----------------------------
// Test setup
// -----------------------------

func setupTestServer() *Server {

	gin.SetMode(gin.TestMode)

	s := &Server{}

	s.Router = gin.Default()

	// 只註冊需要測試的 route
	s.Router.GET("/health", s.HealthHandler)

	return s
}


// -----------------------------
// Health
// -----------------------------

func TestHealthHandler(t *testing.T) {

	srv := setupTestServer()

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	rec := httptest.NewRecorder()

	srv.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	expected := `{"status":"ok"}`
	if rec.Body.String() != expected {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
