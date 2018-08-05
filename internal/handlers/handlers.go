// Package handlers has the main app router and root handler that passes off to other
// handlers, as well as all the web UI handlers and anything else generic
package handlers

import (
	"log"
	"net/http"
)

// RootHandler handles requests for everything, and then compares the requested URL
// to our array of routes, the first match wins and we call that handler.  Requests
// for / are shown a container list, anything not found is 404'd
func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.URL.RawQuery, r.RemoteAddr)
	handler := GetRouteHandler(r.URL.Path)
	if handler != nil {
		handler(w, r)
		return
	}

	ErrorHandler(w, r, 404, "Page not found, double check the URL")
}
