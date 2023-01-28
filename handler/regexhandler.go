// Package hanlder contains the implementation of a simple regex handler
package handler

import (
	"net/http"
	"regexp"
	"strings"
)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexHandler struct {
	routes   []*route
	basePath string
}

var _ http.Handler = &RegexHandler{}

func (h *RegexHandler) SetBasePath(path string) {
	h.basePath = strings.TrimSuffix(path, "/")
}

func (h *RegexHandler) Handle(
	pattern string,
	handler http.Handler,
) {
	h.routes = append(h.routes, &route{
		pattern: regexp.MustCompile(pattern),
		handler: http.StripPrefix(h.basePath, handler),
	})
}

func (h *RegexHandler) HandleFunc(
	pattern string,
	handler func(w http.ResponseWriter, r *http.Request),
) {
	h.routes = append(h.routes, &route{
		pattern: regexp.MustCompile(pattern),
		handler: http.StripPrefix(h.basePath, http.HandlerFunc(handler)),
	})
}

func (h *RegexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
