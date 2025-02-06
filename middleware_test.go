package router_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BambooRaptor/pipeline"
	"github.com/BambooRaptor/router"
)

func addNumToResponse(num int) pipeline.Pipe[http.Handler] {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%v", num)
			next.ServeHTTP(w, r)
		})
	}
}

func TestSeperatePipelinesForRouterAndRoute(t *testing.T) {
	rtr := router.New()
	rtr.Use(
		addNumToResponse(1),
		addNumToResponse(2),
		addNumToResponse(3),
	)

	rootRoute := rtr.Route("/")
	rootRoute.Use(
		addNumToResponse(4),
		addNumToResponse(5),
	).Get(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("\nroot route fin."))
	})

	nested := rootRoute.Route("/nested").Use(
		addNumToResponse(6),
		addNumToResponse(7),
	)
	nested.Get(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("\nnested fin."))
	})

	rootRoute.Route("/other").Use(
		addNumToResponse(8),
		addNumToResponse(9),
	).Get(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("\nother fin."))
	})

	rtr.Route("/nested/deeply/torouter").Use(
		addNumToResponse(10),
		addNumToResponse(11),
	).Get(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("\ndeep fin."))
	})

	nested.Route("/deeply").Use(
		addNumToResponse(12),
		addNumToResponse(13),
	).Get(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("\nnested deep fin."))
	})

	t.Run("Middlware Pipeline", func(t *testing.T) {
		s := httptest.NewServer(rtr)
		defer s.Close()
		assertResponse(t, s, "/", "12345\nroot route fin.")
		assertResponse(t, s, "/nested", "1234567\nnested fin.")
		assertResponse(t, s, "/other", "1234589\nother fin.")
		assertResponse(t, s, "/nested/deeply", "12345671213\nnested deep fin.")
		assertResponse(t, s, "/nested/deeply/torouter", "1231011\ndeep fin.")
	})
}
