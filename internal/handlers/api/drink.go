package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/neophenix/drinkmachine/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Drink is the JSON layout for a drink
type Drink struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Notes       string             `json:"notes"`
	Ingredients []*DrinkIngredient `json:"ingredients"`
}

// DrinkIngredient represents an ingredient in a drink
type DrinkIngredient struct {
	Ingredient string  `json:"ingredient"`
	Amount     float32 `json:"amount"`
	Units      string  `json:"units"`
	Dispense   bool    `json:"dispense"`
	DrinkID    int64   `json:"drink_id"`
}

// DrinkHandler hands off the request to other handlers depending on method and url
func DrinkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// path should be something like /api/drink(/X)?
	parts := strings.Split(r.URL.Path, "/")
	var id int
	if len(parts) == 4 {
		var err error
		id, err = strconv.Atoi(parts[3])
		if err != nil {
			log.Printf("Could not convert drink id %v to int", parts[3])
			JSONError(w, errors.New("Invalid drink identifier"), 400)
			return
		}
	}

	// All the handling for ingredients comes here, so if they wanted a list, pass it off
	if r.Method == "GET" && r.URL.Path == "/api/drink" {
		DrinkListHandler(w, r)
		return
	}

	// Make sure we have a name for all the methods that need it
	if r.Method != "POST" && id == 0 {
		JSONError(w, errors.New("Missing drink identifier"), 400)
		return
	}

	switch r.Method {
	case "GET":
		getDrink(int64(id), w)
	case "POST":
		createDrink(w, r)
	case "PUT":
		updateDrink(int64(id), w, r)
	case "DELETE":
		deleteDrink(int64(id), w)
	}
}

// getDrink returns a single drink to the user
func getDrink(id int64, w http.ResponseWriter) {
	drink, err := models.GetDrink(id, "")
	if err != nil {
		log.Printf("Could not get drink %v: %v\n", id, err.Error())
		JSONError(w, err, 500)
		return
	}

	json, _ := json.Marshal(drink)
	fmt.Fprintf(w, string(json))
}

// createDrink creates a drink and its ingredients
func createDrink(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	buf.ReadFrom(r.Body)

	var data Drink
	err := json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		log.Printf("Error reading incoming data: %v\n", err.Error())
		JSONError(w, err, 400)
		return
	}

	d, err := models.CreateDrink(data.Name, data.Notes)
	if err != nil {
		log.Printf("Error creating drink: %v\n", err.Error())
		JSONError(w, err, 500)
		return
	}

	err = updateDrinkIngredients(d, data.Ingredients)
	if err != nil {
		log.Printf("Error updating drink ingredients: %v\n", err.Error())
		JSONError(w, err, 500)
		return
	}

	// set the ID from the newly made drink so the UI can have it
	data.ID = d.ID
	json, _ := json.Marshal(data)
	fmt.Fprintf(w, string(json))
}

// updateDrink updates a drink and its ingredients
func updateDrink(id int64, w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	buf.ReadFrom(r.Body)

	var data Drink
	err := json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		log.Printf("Error reading incoming data: %v\n", err.Error())
		JSONError(w, err, 400)
		return
	}

	d, err := models.UpdateDrink(id, data.Name, data.Notes)
	if err != nil {
		log.Printf("Error updating drink: %v\n", err.Error())
		JSONError(w, err, 500)
		return
	}

	err = updateDrinkIngredients(d, data.Ingredients)
	if err != nil {
		log.Printf("Error updating drink ingredients: %v\n", err.Error())
		JSONError(w, err, 500)
		return
	}

	json, _ := json.Marshal(data)
	fmt.Fprintf(w, string(json))
}

// updateDrinkIngredients is used by createDrink and updateDrink since they both do roughly the same thing
// it clears out the ingredients first, then adds all the passed ones back in
func updateDrinkIngredients(d *models.Drink, ingredients []*DrinkIngredient) error {
	err := d.ClearIngredients()
	if err != nil {
		return err
	}

	for _, i := range ingredients {
		err = d.AddIngredient(i.Ingredient, i.Amount, i.Units, i.Dispense)
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteDrink removes a drink and its ingredients
func deleteDrink(id int64, w http.ResponseWriter) {
	err := models.DeleteDrink(id)
	if err != nil {
		log.Printf("Error deleting drink %v: %v\n", id, err.Error())
		JSONError(w, err, 500)
		return
	}

	fmt.Fprintf(w, fmt.Sprintf("{\"id\":%v}", id))
}

// DrinkListHandler will return all the drinks in the system, adding ?format=simple to
// the URL will cause it to not populate ingredients, which isn't needed for a list page
func DrinkListHandler(w http.ResponseWriter, r *http.Request) {
	simple := false
	// keep this easy until I realize I need more query params
	if r.URL.RawQuery == "format=simple" {
		simple = true
	}

	list, err := models.GetAllDrinks(simple)
	if err != nil {
		log.Printf("Could not get list of drinks: %v\n", err.Error())
		JSONError(w, err, 500)
		return
	}

	json, _ := json.Marshal(list)
	fmt.Fprintf(w, string(json))
}
