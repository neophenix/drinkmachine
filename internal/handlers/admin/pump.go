package admin

import (
	"bytes"
	"fmt"
	"github.com/neophenix/drinkmachine/internal/hw"
	"github.com/neophenix/drinkmachine/internal/models"
	"github.com/neophenix/drinkmachine/internal/template"
	"log"
	"net/http"
)

// PumpHandler serves the "Manage Pumps" page
func PumpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	ingredients, err := models.GetAllIngredients()
	if err != nil {
		msg := fmt.Sprintf("Could not get ingredient list: %v\n", err.Error())
		log.Printf(msg)
		ErrorHandler(w, r, 500, msg)
		return
	}

	tmpl := template.ReadTemplate("admin/pumps.tmpl")

	var out bytes.Buffer
	tmpl.ExecuteTemplate(&out, "base", map[string]interface{}{
		"Pumps":       hw.Pumps,
		"Ingredients": ingredients,
	})

	fmt.Fprintf(w, string(out.Bytes()))
}
