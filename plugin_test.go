package traefik_plugin_rename_header_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	plugin "github.com/beyerleinf/traefik-plugin-rename-header"
)

func TestCreateConfig(t *testing.T) {
	config := plugin.CreateConfig()
	if config == nil {
		t.Fatal("CreateConfig() returned nil")
	}
	if config.OldHeader != "" {
		t.Errorf("Expected empty OldHeader, got %q", config.OldHeader)
	}
	if config.NewHeader != "" {
		t.Errorf("Expected empty NewHeader, got %q", config.NewHeader)
	}
}

func TestNewValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *plugin.Config
		expectError bool
	}{
		{
			name: "valid configuration",
			config: &plugin.Config{
				OldHeader: "X-Old",
				NewHeader: "X-New",
			},
			expectError: false,
		},
		{
			name: "missing old header",
			config: &plugin.Config{
				NewHeader: "X-New",
			},
			expectError: true,
		},
		{
			name: "missing new header",
			config: &plugin.Config{
				OldHeader: "X-Old",
			},
			expectError: true,
		},
		{
			name:        "empty config",
			config:      &plugin.Config{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
			_, err := plugin.New(context.Background(), next, tt.config, "test")

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestPlugin(t *testing.T) {
	tests := []struct {
		name            string
		oldHeader       string
		newHeader       string
		inputHeaders    map[string]string
		expectedHeaders map[string]string
	}{
		{
			name:      "basic header rename",
			oldHeader: "X-Old",
			newHeader: "X-New",
			inputHeaders: map[string]string{
				"X-Old": "test-value",
			},
			expectedHeaders: map[string]string{
				"X-New": "test-value",
			},
		},
		{
			name:      "no header present",
			oldHeader: "X-Old",
			newHeader: "X-New",
			inputHeaders: map[string]string{
				"X-Other": "other-value",
			},
			expectedHeaders: map[string]string{
				"X-Other": "other-value",
			},
		},
		{
			name:      "multiple headers",
			oldHeader: "X-Old",
			newHeader: "X-New",
			inputHeaders: map[string]string{
				"X-Old":   "test-value",
				"X-Other": "other-value",
			},
			expectedHeaders: map[string]string{
				"X-New":   "test-value",
				"X-Other": "other-value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedHeaders http.Header
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				receivedHeaders = req.Header.Clone()
			})

			config := plugin.CreateConfig()
			config.OldHeader = tt.oldHeader
			config.NewHeader = tt.newHeader

			ctx := context.Background()

			handler, err := plugin.New(ctx, next, config, "plugin")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			if err != nil {
				t.Fatal(err)
			}

			for k, v := range tt.inputHeaders {
				req.Header.Set(k, v)
			}

			handler.ServeHTTP(recorder, req)

			for k, v := range tt.expectedHeaders {
				if got := receivedHeaders.Get(k); got != v {
					t.Errorf("Expected header %q to be %q, got %q", k, v, got)
				}
			}

			if tt.inputHeaders[tt.oldHeader] != "" {
				if got := receivedHeaders.Get(tt.oldHeader); got != "" {
					t.Errorf("Expected old header %q to be removed, but got value %q", tt.oldHeader, got)
				}
			}
		})
	}
}

func TestMiddlewareChain(t *testing.T) {
	handlerCalled := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handlerCalled = true
		rw.WriteHeader(http.StatusOK)
	})

	config := plugin.CreateConfig()
	config.OldHeader = "X-Old"
	config.NewHeader = "X-New"

	handler, err := plugin.New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatalf("Failed to create middleware: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	rw := httptest.NewRecorder()

	handler.ServeHTTP(rw, req)

	if !handlerCalled {
		t.Error("Next handler was not called")
	}
	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}
}
