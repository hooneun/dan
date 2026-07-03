package dan

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoveryHandlesRuntimeFailure(t *testing.T) {
	app := NewEngine()

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

func TestBodyLimit(t *testing.T) {
	app := NewEngine()
	app.Use(BodyLimit(4))

	app.POST("/body", func(c *Context) error {
		body, err := io.ReadAll(c.R.Body)
		if err != nil {
			return c.Error(http.StatusRequestEntityTooLarge, "Request body too large")
		}

		return c.JSON(http.StatusOK, map[string]string{"body": string(body)})
	})

	req := httptest.NewRequest(http.MethodPost, "/body", strings.NewReader("too-large"))
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d", rec.Code)
	}
}

func TestSanitizeLogValue(t *testing.T) {
	got := sanitizeLogValue("/users\nforged\rline\ttab")
	want := `/users\nforged\rline\ttab`

	if got != want {
		t.Fatalf("expected sanitized value %q, got %q", want, got)
	}
}
