package handlers

import (
	"bytes"
	"fmt"
	"github.com/neophenix/drinkmachine/internal/template"
	"net/http"
)

// ErrorHandler is a generic way to report an error to the user and set the HTTP status code
func ErrorHandler(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	w.Header().Set("Content-Type", "text/html")

	tmpl := template.ReadTemplate("error.tmpl")

	var out bytes.Buffer
	tmpl.ExecuteTemplate(&out, "base", map[string]interface{}{
		"Message": msg,
	})

	w.WriteHeader(statusCode)
	fmt.Fprintf(w, string(out.Bytes()))
}
