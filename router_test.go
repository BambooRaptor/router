package router_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BambooRaptor/router"
	"github.com/BambooRaptor/router/pkgs/set"
)

func TestRouter(t *testing.T) {
	r := router.New()

	r.Route("/").GetFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			t.Fatalf("Failed to write HTTP response: %v", err)
		}
	})

	assertResponse := func(t *testing.T, s *httptest.Server, expected string) {
		resp, err := s.Client().Get(s.URL)
		if err != nil {
			t.Fatalf("Unexpected error from server: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Fatalf("Status Code expected, but got:\n[200] <=/=> [%v]", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error reading body: %v", err)
		}

		if !bytes.Equal(body, []byte(expected)) {
			t.Fatalf("Response expected, but got:\n%q <=/=> %q", body, expected)
		}
	}

	t.Run("basic router", func(t *testing.T) {
		s := httptest.NewServer(r)
		defer s.Close()
		assertResponse(t, s, "Hello, World!")
	})

	t.Run("basic TLSS router", func(t *testing.T) {
		s := httptest.NewTLSServer(r)
		defer s.Close()
		assertResponse(t, s, "Hello, World!")
	})
}

func TestAllowedMethods(t *testing.T) {
	r := router.New()
	r.Use(r.SetAllowedMethods)

	r.Route("/ping").GetFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("get-pong"))
		if err != nil {
			t.Fatalf("Failed to write HTTP response: %v", err)
		}
	})

	r.Route("/ping").PostFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("post-pong"))
		if err != nil {
			t.Fatalf("Failed to write HTTP response: %v", err)
		}
	})

	t.Run("Set Access Methods", func(t *testing.T) {
		s := httptest.NewServer(r)
		defer s.Close()
		resp, err := s.Client().Get(s.URL + "/ping")
		if err != nil {
			t.Fatalf("Unexpected error from server: %v", err)
		}

		header := resp.Header.Get("Access-Control-Allow-Methods")

		methods := set.FromArray(strings.Split(header, ", "))
		expected := set.FromArray([]string{"GET", "POST"})

		if !expected.Matches(methods) {
			t.Fatalf("Response expected, but got:\n%q <=/=> %q", expected, methods)
		}
	})
}
