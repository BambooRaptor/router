package router

import (
	"net/http"

	"github.com/BambooRaptor/pipeline"
)

type Router struct {
	mux    *http.ServeMux
	routes map[string]*route
	pipe   pipeline.Pipeline[http.Handler]
}

// Return a new router that can be built upon.
// Route -> Attach a new route to the router
// Use -> Define middleware to use with the route
// UsePipeline -> Define a middleware pipeline to use
func New() *Router {
	mux := http.NewServeMux()
	return &Router{mux, make(map[string]*route), pipeline.New[http.Handler]()}
}

// Creates a new, validated route
// OR
// Returns a route if it already exists
func (rtr *Router) Route(path string) *route {
	nrt := newRoute(rtr, path)
	rt, exists := rtr.routes[path]
	if exists {
		return rt
	}
	rtr.routes[path] = nrt
	return nrt
}

// Attach middleware globally to the router
// This is used for any route attached to the router
func (rtr *Router) Use(funcs ...pipeline.Pipe[http.Handler]) {
	rtr.pipe = rtr.pipe.IntoRaw(funcs...)
}

// Attach middleware pipeline globally to the router
// This is used for any route attached to the router
func (rtr *Router) UsePipeline(pipe pipeline.Pipeline[http.Handler]) {
	rtr.pipe = rtr.pipe.Into(pipe)
}

// [http.Handler] interface implementation
func (rout *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rout.mux.ServeHTTP(w, r)
}

// DEBUG
func (rtr *Router) GetAllRoutes() []*route {
	routes := make([]*route, 0)
	for _, rout := range rtr.routes {
		routes = append(routes, rout)
	}
	return routes
}
