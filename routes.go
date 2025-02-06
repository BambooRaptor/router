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

// Create a new route using this route as a base.
// Also copies over the middleware pipeline to the new route
func (r *route) Route(pattern string) *route {
	return r.router.Route(r.path + pattern).UsePipeline(r.pipe)
}

// Debugging purposes -> Returns the path of the route
func (r *route) String() string {
	return r.path
}

// Attach middleware specifically to this route,
// and subroutes derived from this route
func (r *route) Use(funcs ...pipeline.Pipe[http.Handler]) *route {
	r.pipe = r.pipe.IntoRaw(funcs...)
	return r
}

// Attach middleware pipeline specifically to this route,
// and subroutes derived from this route
func (r *route) UsePipeline(pipe pipeline.Pipeline[http.Handler]) *route {
	r.pipe = r.pipe.Into(pipe)
	return r
}

// Handler the route with a custom method and handler
func (r *route) Handler(method string, handler http.Handler) {
	path := r.path

	if method != "" {
		if r.methods.Has(method) {
			panic(fmt.Sprintf("Method [%s] on route %q already exists", method, path))
		}
		err := r.methods.Add(method)
		if err != nil {
			panic(fmt.Sprintf("Route %q with method [%v] already exists", path, method))
		}
		path = method + " " + r.path
	}

	r.router.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		r.router.pipe.Into(r.pipe).Build(handler).ServeHTTP(w, req)
	})
}

// Handle the route with a custom method and function
func (r *route) Handle(method string, handler http.HandlerFunc) {
	r.Handler(method, handler)
}

// Get all the methods currently implemented on the route
func (r *route) GetMethods() []string {
	return r.methods.ToArray()
}

// Handle the route with the GET method and a handler
func (r *route) GetHandler(handler http.Handler) {
	r.Handler("GET", handler)
}

// Handle the route with the GET method and a function
func (r *route) Get(handler http.HandlerFunc) {
	r.GetHandler(handler)
}

// Handle the route with the POST method and a handler
func (r *route) PostHandler(handler http.Handler) {
	r.Handler("POST", handler)
}

// Handle the route with the POST method and a function
func (r *route) Post(handler http.HandlerFunc) {
	r.PostHandler(handler)
}

// Handle the route with the PUT method and a handler
func (r *route) PutHandler(handler http.Handler) {
	r.Handler("PUT", handler)
}

// Handle the route with the PUT method and a function
func (r *route) Put(handler http.HandlerFunc) {
	r.PutHandler(handler)
}

// Handle the route with the DELETE method and a handler
func (r *route) DeleteHandler(handler http.Handler) {
	r.Handler("DELETE", handler)
}

// Handle the route with the DELETE method and a function
func (r *route) Delete(handler http.HandlerFunc) {
	r.DeleteHandler(handler)
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
