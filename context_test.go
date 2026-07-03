package dan

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/search?q=dan&page=1", nil)
	ctx := &Context{W: httptest.NewRecorder(), R: req}

	if got := ctx.Query("q"); got != "dan" {
		t.Fatalf("expected query q to be dan, got %q", got)
	}

	if got := ctx.Query("missing"); got != "" {
		t.Fatalf("expected missing query to be empty, got %q", got)
	}
}

func TestDefaultQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/search?q=dan", nil)
	ctx := &Context{W: httptest.NewRecorder(), R: req}

	if got := ctx.DefaultQuery("q", "default"); got != "dan" {
		t.Fatalf("expected query q to be dan, got %q", got)
	}

	if got := ctx.DefaultQuery("page", "1"); got != "1" {
		t.Fatalf("expected default query page to be 1, got %q", got)
	}
}

func TestForm(t *testing.T) {
	body := strings.NewReader("name=dan&role=framework")
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := &Context{W: httptest.NewRecorder(), R: req}

	if got := ctx.Form("name"); got != "dan" {
		t.Fatalf("expected form name to be dan, got %q", got)
	}

	if got := ctx.Form("missing"); got != "" {
		t.Fatalf("expected missing form to be empty, got %q", got)
	}
}

func TestDefaultForm(t *testing.T) {
	body := strings.NewReader("name=dan")
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := &Context{W: httptest.NewRecorder(), R: req}

	if got := ctx.DefaultForm("name", "default"); got != "dan" {
		t.Fatalf("expected form name to be dan, got %q", got)
	}

	if got := ctx.DefaultForm("role", "framework"); got != "framework" {
		t.Fatalf("expected default form role to be framework, got %q", got)
	}
}
