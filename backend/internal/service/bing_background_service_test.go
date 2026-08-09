package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBingBackgroundServiceResolve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"images":[{"url":"/th?id=example"}]}`))
	}))
	defer server.Close()
	service := newBingBackgroundService(server.Client(), server.URL)
	got, err := service.Resolve(context.Background())
	if err != nil || got != "https://www.bing.com/th?id=example" {
		t.Fatalf("Resolve() = %q, %v", got, err)
	}
}

func TestBingBackgroundServiceRejectsExternalURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"images":[{"url":"https://example.com/image.jpg"}]}`))
	}))
	defer server.Close()
	_, err := newBingBackgroundService(server.Client(), server.URL).Resolve(context.Background())
	if err == nil || !strings.HasPrefix(err.Error(), "ERR-1011:") {
		t.Fatalf("expected ERR-1011, got %v", err)
	}
}
