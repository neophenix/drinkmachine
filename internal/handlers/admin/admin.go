package admin

import (
	"bytes"
	"fmt"
	"github.com/neophenix/drinkmachine/internal/template"
	"net/http"
)

// Handler serves the /admin page, if we decide to make that
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl := template.ReadTemplate("admin/admin.tmpl")

	var out bytes.Buffer
	tmpl.ExecuteTemplate(&out, "base", map[string]interface{}{})

	fmt.Fprintf(w, string(out.Bytes()))
}

// ErrorHandler is a generic way to report an error to the user and set the HTTP status code
// copied from handlers ... I thought I might want to do something different in the admin section but so far no
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
