package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	
)


func TestHealthHandler(t *testing.T){

	srv := setupTestServer()

	srv.Router.GET(
		"/health",
		srv.HealthHandler,
	)


	req :=
		httptest.NewRequest(
			http.MethodGet,
			"/health",
			nil,
		)


	rec :=
		httptest.NewRecorder()


	srv.Router.ServeHTTP(
		rec,
		req,
	)


	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}

}
