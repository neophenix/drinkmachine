package handlers

import (
	"bytes"
	"fmt"
	"github.com/neophenix/drinkmachine/internal/template"
	"net/http"
)

// FourOhFourHandler is our 404 response
func FourOhFourHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl := template.ReadTemplate("404.tmpl")

	var out bytes.Buffer
	tmpl.ExecuteTemplate(&out, "base", map[string]interface{}{
		"Page": "none",
	})

	fmt.Fprintf(w, string(out.Bytes()))
}
