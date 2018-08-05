package handlers

import (
	"bytes"
	"fmt"
	"github.com/neophenix/drinkmachine/internal/models"
	"github.com/neophenix/drinkmachine/internal/template"
	"log"
	"net/http"
)

// HomeHandler takes care of really the only non admin page, / or "pour me a drink"
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	drinks, err := models.GetAllDrinks(true)
	if err != nil {
		msg := fmt.Sprintf("Could not get drink list: %v\n", err.Error())
		log.Printf(msg)
		ErrorHandler(w, r, 500, msg)
		return
	}

	tmpl := template.ReadTemplate("home.tmpl")

	var out bytes.Buffer
	tmpl.ExecuteTemplate(&out, "base", map[string]interface{}{
		"Drinks": drinks,
	})

	fmt.Fprintf(w, string(out.Bytes()))
}
