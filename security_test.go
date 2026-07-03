package dan

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerErrorDoesNotOverwriteWrittenResponse(t *testing.T) {
	app := NewEngine()

	app.GET("/late-error", func(c *Context) error {
		if err := c.JSON(http.StatusOK, map[string]string{"status": "ok"}); err != nil {
			return err
		}

		return errors.New("late error")
	})

	req := httptest.NewRequest(http.MethodGet, "/late-error", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `"status":"ok"`) {
		t.Fatalf("expected original response body, got %s", body)
	}

	if strings.Contains(body, "Internal Server Error") {
		t.Fatalf("expected no duplicate error response, got %s", body)
	}
}
