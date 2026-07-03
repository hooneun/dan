package dan

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoveryHandlesRuntimeFailure(t *testing.T) {
	app := NewEngine()
	app.Use(Recovery())

	app.GET("/fail", func(c *Context) error {
		panic(errors.New("handler failed"))
	})

	req := httptest.NewRequest(http.MethodGet, "/fail", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"error":"Internal Server Error"`) {
		t.Fatalf("expected internal server error response, got %s", rec.Body.String())
	}
}

func TestRecoveryPassesThroughHandler(t *testing.T) {
	app := NewEngine()
	app.Use(Recovery())

	app.GET("/ok", func(c *Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Fatalf("expected ok response, got %s", rec.Body.String())
	}
}
