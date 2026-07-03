package dan

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDynamicRoute(t *testing.T) {
	app := NewEngine()

	app.GET("/users/:id", func(c *Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"id": c.Param("id"),
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"id":"123"`) {
		t.Fatalf("expected id param, got %s", rec.Body.String())
	}
}

func TestGroupedDynamicRoute(t *testing.T) {
	app := NewEngine()
	api := app.Group("/api/v1")

	api.GET("/users/:id", func(c *Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"id": c.Param("id"),
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/456", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"id":"456"`) {
		t.Fatalf("expected id param, got %s", rec.Body.String())
	}
}

func TestStaticRouteHasPriority(t *testing.T) {
	app := NewEngine()

	app.GET("/users/:id", func(c *Context) error {
		return c.JSON(http.StatusOK, map[string]string{"route": "dynamic"})
	})

	app.GET("/users/me", func(c *Context) error {
		return c.JSON(http.StatusOK, map[string]string{"route": "static"})
	})

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"route":"static"`) {
		t.Fatalf("expected static route, got %s", rec.Body.String())
	}
}

func TestDynamicRouteNotFoundWhenSegmentCountDiffers(t *testing.T) {
	app := NewEngine()

	app.GET("/users/:id", func(c *Context) error {
		return c.JSON(http.StatusOK, nil)
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123/posts", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
