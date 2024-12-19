package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"jst.dev/src/web"
)

func TestHandler(t *testing.T) {
	navItemsIndex := []web.NavItem{
		{Href: "/", Text: "Home", Active: true},
	}
	navItemsNotFound := []web.NavItem{
		{Href: "/404", Text: "404", Active: false},
	}
	s := &Server{}
	server := httptest.NewServer(s.handlerRoot(navItemsIndex, navItemsNotFound))
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expectedStrings := []string{
		">Home</h1>",
		">Work in progress..</p>",
		"/assets/js/page/index.js",
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(string(body), expected) {
			t.Errorf("expected response body to contain %v; got %v", expected, string(body))
		}
	}
}
