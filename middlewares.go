package router

import (
	"net/http"
	"strings"
)

func (router *Router) SetAllowedMethods(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pattern := r.Pattern
		patternHasMethod := strings.ContainsAny(r.Pattern, " ")
		if patternHasMethod {
			pattern = strings.Split(r.Pattern, " ")[1]
		}

		route, ok := router.routes[pattern]
		if ok {
			if len(route.methods) > 0 {
				w.Header().Set(
					"Access-Control-Allow-Methods",
					strings.Join(route.GetMethods(),
						", ",
					))
			} else {
				w.Header().Set("Access-Control-Allow-Methods", "*")
			}
		}
		next.ServeHTTP(w, r)
	})
}
