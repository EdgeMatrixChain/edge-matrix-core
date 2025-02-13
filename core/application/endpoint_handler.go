package application

import (
	"net/http"
)

type EndpointHandler struct {
	routes map[string]func(w http.ResponseWriter, r *http.Request)
}

func (h *EndpointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := h.routes[r.URL.Path]; ok {
		handler(w, r)
	} else {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}

func (h *EndpointHandler) AddHandler(url string, handler func(w http.ResponseWriter, r *http.Request)) {
	h.routes[url] = handler
}
