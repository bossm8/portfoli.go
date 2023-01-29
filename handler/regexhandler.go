// Copyright (c) 2023, Boss Marco <bossm8@hotmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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

// RegexHandler implements a http handler with regex pattern matching
type RegexHandler struct {
	routes   []*route
	basePath string
}

// Make sure the Handler interface is implemented
var _ http.Handler = &RegexHandler{}

// SetBasePath sets the base path of the server to path
// This path will be stripped before the request is passed to any handler
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
