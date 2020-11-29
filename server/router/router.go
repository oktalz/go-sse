package router

import (
	"log"
	"net/http"
)

// Router serves http
type Router struct {
	handlers    map[string]func(http.ResponseWriter, *http.Request)
	middlewares []func(http.Handler) http.Handler
}

func New() *Router {
	return &Router{
		handlers: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
}

func (rt *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	rt.middlewares = append(rt.middlewares, middlewares...)
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	log.Println(url)
	f, ok := rt.handlers[""]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	f(w, r)
}
