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

func TestServer_handlerRedirectWithUrlParam(t *testing.T) {
	type args struct {
		url       string
		paramKeys []string
		status    int
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
		setup     func(*httptest.ResponseRecorder, *http.Request)
		validate  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful redirect with single param",
			args: args{
				url:       "/url/i/{shortCode}",
				paramKeys: []string{"shortCode"},
				status:    http.StatusMovedPermanently,
			},
			wantPanic: false,
			setup: func(w *httptest.ResponseRecorder, r *http.Request) {
				r.SetPathValue("shortCode", "abc123")
			},
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				if w.Code != http.StatusMovedPermanently {
					t.Errorf("expected status %d; got %d", http.StatusMovedPermanently, w.Code)
				}
				if loc := w.Header().Get("Location"); loc != "/url/i/abc123" {
					t.Errorf("expected redirect to %s; got %s", "/url/i/abc123", loc)
				}
			},
		},
		{
			name: "panic when param key not in url",
			args: args{
				url:       "/url/i/static",
				paramKeys: []string{"shortCode"},
				status:    http.StatusMovedPermanently,
			},
			wantPanic: true,
		},
		{
			name: "successful redirect with multiple params",
			args: args{
				url:       "/url/{category}/{id}",
				paramKeys: []string{"category", "id"},
				status:    http.StatusTemporaryRedirect,
			},
			wantPanic: false,
			setup: func(w *httptest.ResponseRecorder, r *http.Request) {
				r.SetPathValue("category", "blog")
				r.SetPathValue("id", "123")
			},
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				if w.Code != http.StatusTemporaryRedirect {
					t.Errorf("expected status %d; got %d", http.StatusTemporaryRedirect, w.Code)
				}
				if loc := w.Header().Get("Location"); loc != "/url/blog/123" {
					t.Errorf("expected redirect to %s; got %s", "/url/blog/123", loc)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("handlerRedirectWithUrlParam() panic = %v, wantPanic = %v", r != nil, tt.wantPanic)
				}
			}()

			s := &Server{}
			handler := s.handlerRedirectWithUrlParam(tt.args.url, tt.args.paramKeys, tt.args.status)

			if !tt.wantPanic {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/test", nil)
				if tt.setup != nil {
					tt.setup(w, r)
				}
				handler.ServeHTTP(w, r)
				if tt.validate != nil {
					tt.validate(t, w)
				}
			}
		})
	}
}
