package router

import (
	"github.com/BambooRaptor/pipeline"
	"net/http"
	"strings"
)

type route struct {
	*http.ServeMux
	base string
	pipe pipeline.Pipeline[http.Handler]
}

func newRoute(base string, pipe pipeline.Pipeline[http.Handler]) *route {
	return &route{nil, base, pipe}
}

func (r *route) AssignMux(mux *http.ServeMux) {
	r.ServeMux = mux
}

func (r *route) Route(pattern string) *route {
	pattern = sanitizeRoute(pattern)

	if len(pattern) == 0 {
		panic("Cannot have an empty route")
	}

	if pattern[len(pattern)-1] != '/' {
		panic("Routes cannot end with '/'")
	}

	if pattern[0] != '/' {
		panic("Routes must begin with '/'")
	}

	return &route{
		ServeMux: r.ServeMux,
		base:     r.base + pattern,
		pipe:     r.pipe,
	}
}

func (r *route) Use(funcs ...pipeline.Pipe[http.Handler]) *route {
	r.pipe = r.pipe.IntoRaw(funcs...)
	return r
}

func (r *route) UsePipeline(pipe pipeline.Pipeline[http.Handler]) *route {
	r.pipe = r.pipe.Into(pipe)
	return r
}

func (r *route) Handle(method string, handler http.Handler) {
	r.base = sanitizeRoute(r.base)
	path := r.base
	if method != "" {
		path = method + " " + r.base
	}
	r.ServeMux.Handle(path, r.pipe.Build(handler))
}

func (r *route) HandleFunc(method string, handler http.HandlerFunc) {
	r.Handle(method, handler)
}

// Get
func (r *route) Get(handler http.Handler) {
	r.Handle("GET", handler)
}

func (r *route) GetFunc(handler http.HandlerFunc) {
	r.Get(handler)
}

// Post
func (r *route) Post(handler http.Handler) {
	r.Handle("POST", handler)
}

func (r *route) PostFunc(handler http.HandlerFunc) {
	r.Post(handler)
}

func (r *route) Put(handler http.Handler) {
	r.Handle("PUT", handler)
}

func (r *route) PutFunc(handler http.HandlerFunc) {
	r.Put(handler)
}

func (r *route) Delete(handler http.Handler) {
	r.Handle("DELETE", handler)
}

func (r *route) DeleteFunc(handler http.HandlerFunc) {
	r.Delete(handler)
}

// UTIL FUNCS
func sanitizeRoute(route string) string {
	for strings.Contains(route, "//") {
		route = strings.ReplaceAll(route, "//", "/")
	}

	if route[len(route)-1] != '/' {
		route += "/"
	}

	return route
}
