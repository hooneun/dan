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

func TestMoreHTTPMethods(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{name: "PUT", method: http.MethodPut},
		{name: "PATCH", method: http.MethodPatch},
		{name: "DELETE", method: http.MethodDelete},
		{name: "OPTIONS", method: http.MethodOptions},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewEngine()

			switch tt.method {
			case http.MethodPut:
				app.PUT("/resources/:id", methodHandler(tt.method))
			case http.MethodPatch:
				app.PATCH("/resources/:id", methodHandler(tt.method))
			case http.MethodDelete:
				app.DELETE("/resources/:id", methodHandler(tt.method))
			case http.MethodOptions:
				app.OPTIONS("/resources/:id", methodHandler(tt.method))
			}

			req := httptest.NewRequest(tt.method, "/resources/123", nil)
			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", rec.Code)
			}

			body := rec.Body.String()
			if !strings.Contains(body, `"method":"`+tt.method+`"`) {
				t.Fatalf("expected method %s, got %s", tt.method, body)
			}

			if !strings.Contains(body, `"id":"123"`) {
				t.Fatalf("expected id param, got %s", body)
			}
		})
	}
}

func TestMethodSpecificRoutes(t *testing.T) {
	app := NewEngine()

	app.GET("/resources/:id", methodHandler(http.MethodGet))
	app.PUT("/resources/:id", methodHandler(http.MethodPut))

	req := httptest.NewRequest(http.MethodPut, "/resources/123", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"method":"PUT"`) {
		t.Fatalf("expected PUT handler, got %s", rec.Body.String())
	}
}

func methodHandler(method string) HandlerFunc {
	return func(c *Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"id":     c.Param("id"),
			"method": method,
		})
	}
}
