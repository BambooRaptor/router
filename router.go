package router

import (
	"net/http"

	"github.com/BambooRaptor/pipeline"
	"github.com/BambooRaptor/router/pkgs/set"
)

type Router struct {
	mux    *http.ServeMux
	routes map[string]*route
	pipe   pipeline.Pipeline[http.Handler]
}

// Return a new route that handles the "/" base case
func New() *Router {
	mux := http.NewServeMux()
	return &Router{mux, make(map[string]*route), pipeline.New[http.Handler]()}
}

func (rout *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rout.mux.ServeHTTP(w, r)
}

func (rtr *Router) Use(funcs ...pipeline.Pipe[http.Handler]) {
	rtr.pipe = rtr.pipe.IntoRaw(funcs...)
}

func (rtr *Router) UsePipeline(pipe pipeline.Pipeline[http.Handler]) {
	rtr.pipe = rtr.pipe.Into(pipe)
}

func (rtr *Router) Route(path string) *route {
	return rtr.newRoute(path, rtr.pipe)
}

func (rtr *Router) newRoute(path string, pipe pipeline.Pipeline[http.Handler]) *route {
	path = sanitizeRoute(path)

	if len(path) == 0 {
		panic("Cannot have an empty route")
	}

	if len(path) > 1 && path[len(path)-1] == '/' {
		panic("Routes cannot end with '/'")
	}

	if path[0] != '/' {
		panic("Routes must begin with '/'")
	}

	rt, ok := rtr.routes[path]

	if ok {
		return rt
	}
	rt = &route{rtr, path, pipe, *set.New[string]()}
	rtr.routes[path] = rt
	return rt
}

func (rtr *Router) GetAllRoutes() []*route {
	routes := make([]*route, 0)
	for _, rout := range rtr.routes {
		routes = append(routes, rout)
	}
	return routes
}
