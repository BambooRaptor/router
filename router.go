package router

import (
	"net/http"

	"github.com/BambooRaptor/pipeline"
)

func NewRouter() *route {
	mux := http.NewServeMux()
	route := newRoute("/", pipeline.Empty[http.Handler]())
	route.AssignMux(mux)
	return route
}
