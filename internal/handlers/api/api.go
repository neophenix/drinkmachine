// Package api has all the API handlers and utilities specific to the api
package api

import (
	"fmt"
	"net/http"
)

// JSONError takes an error and turns it into a JSON string like {"error": ...} and sets the passed HTTP status code
// make the string by "hand" because its simple and easier to do
func JSONError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "{\"error\": \"%s\"}", err.Error())
}
