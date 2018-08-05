package admin

import (
	"bytes"
	"fmt"
	"github.com/neophenix/drinkmachine/internal/models"
	"github.com/neophenix/drinkmachine/internal/template"
	"log"
	"net/http"
)

// DrinkHandler serves the "Manage Drinks" page
func DrinkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	drinks, err := models.GetAllDrinks(true)
	if err != nil {
		msg := fmt.Sprintf("Could not get drink list: %v\n", err.Error())
		log.Printf(msg)
		ErrorHandler(w, r, 500, msg)
		return
	}

	ingredients, err := models.GetAllIngredients()
	if err != nil {
		msg := fmt.Sprintf("Could not get ingredient list: %v\n", err.Error())
		log.Printf(msg)
		ErrorHandler(w, r, 500, msg)
		return
	}

	tmpl := template.ReadTemplate("admin/drinks.tmpl")

	var out bytes.Buffer
	tmpl.ExecuteTemplate(&out, "base", map[string]interface{}{
		"Drinks":      drinks,
		"Ingredients": ingredients,
	})

	fmt.Fprintf(w, string(out.Bytes()))
}
