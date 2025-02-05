package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/BambooRaptor/pipeline"
	"github.com/BambooRaptor/router/pkgs/set"
)

type route struct {
	router  *Router
	path    string
	pipe    pipeline.Pipeline[http.Handler]
	methods set.Set[string]
}

func newRoute(router *Router, path string) *route {
	path = validatePath(path)
	return &route{router, path, pipeline.New[http.Handler](), set.New[string]()}
}

func (r *route) Route(pattern string) *route {
	return r.router.Route(r.path + pattern).UsePipeline(r.pipe)
}

func (r *route) String() string {
	return fmt.Sprint(r.path)
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
	path := r.path

	if method != "" {
		if r.methods.Has(method) {
			panic(fmt.Sprintf("Method [%s] on route %q already exists", method, path))
		}
		r.methods.Add(method)
		path = method + " " + r.path
	}

	r.router.mux.Handle(path, r.router.pipe.Into(r.pipe).Build(handler))
}

func (r *route) HandleFunc(method string, handler http.HandlerFunc) {
	r.Handle(method, handler)
}

func (r *route) GetMethods() []string {
	return r.methods.ToArray()
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

func (r *route) GetAllRoutes() []*route {
	return r.router.GetAllRoutes()
}

// UTIL FUNCS
func sanitizePath(route string) string {
	for strings.Contains(route, "//") {
		route = strings.ReplaceAll(route, "//", "/")
	}

	// if route[len(route)-1] != '/' {
	// 	route += "/"
	// }

	return route
}

func validatePath(path string) string {
	path = sanitizePath(path)
	if len(path) == 0 {
		panic("Cannot have an empty route")
	}

	if len(path) > 1 && path[len(path)-1] == '/' {
		panic("Routes cannot end with '/'")
	}

	if path[0] != '/' {
		panic("Routes must begin with '/'")
	}
	return path
}
