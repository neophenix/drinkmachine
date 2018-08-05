package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/neophenix/drinkmachine/internal/models"
	"log"
	"net/http"
	"strings"
)

// Ingredient is the JSON struct for an ingredient
type Ingredient struct {
	Ingredient string  `json:"ingredient"`
	Viscosity  float32 `json:"viscosity"`
}

// IngredientHandler calls other "handlers" depending on request method
func IngredientHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// path should be something like /api/ingredient/X
	parts := strings.Split(r.URL.Path, "/")
	var name string
	if len(parts) == 4 {
		name = parts[3]
	}

	// All the handling for ingredients comes here, so if they wanted a list, pass it off
	if r.Method == "GET" && r.URL.Path == "/api/ingredient" {
		IngredientListHandler(w, r)
		return
	}

	// Make sure we have a name for all the methods that need it
	if r.Method != "POST" && name == "" {
		JSONError(w, errors.New("Missing ingredient name"), 400)
		return
	}

	// return the ingredient info
	switch r.Method {
	case "GET":
		getIngredient(name, w)
	case "POST":
		createIngredient(w, r)
	case "PUT":
		updateIngredient(name, w, r)
	case "DELETE":
		deleteIngredient(name, w)
	}
}

// getIngredient is called by IngredientHandler for GETs and returns the ingredient info to the user
func getIngredient(name string, w http.ResponseWriter) {
	row, err := models.GetIngredient(name)
	if err != nil {
		log.Printf("Error getting ingredient: %v\n", err.Error())
		JSONError(w, err, 500)
		return
	}
	json, _ := json.Marshal(row)
	fmt.Fprintf(w, string(json))
}

// createIngredient is called by IngredientHandler for POSTs and creates an ingredient
func createIngredient(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	buf.ReadFrom(r.Body)

	var data Ingredient
	err := json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		log.Printf("Error reading incoming data: %v\n", err.Error())
		JSONError(w, err, 400)
		return
	}

	err = models.CreateIngredient(data.Ingredient, data.Viscosity)
	if err != nil {
		log.Printf("Error creating ingredient %v: %v\n", data.Ingredient, err.Error())
		JSONError(w, err, 500)
		return
	}

	json, _ := json.Marshal(data)
	fmt.Fprintf(w, string(json))
}

// updateIngredient is called by IngredientHandler for PUTs and updates an ingredient
func updateIngredient(name string, w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	buf.ReadFrom(r.Body)

	var data Ingredient
	err := json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		log.Printf("Error reading incoming data: %v\n", err.Error())
		JSONError(w, err, 400)
		return
	}

	err = models.UpdateIngredient(name, data.Viscosity)
	if err != nil {
		log.Printf("Error updating ingredient %v: %v\n", name, err.Error())
		JSONError(w, err, 500)
		return
	}

	json, _ := json.Marshal(data)
	fmt.Fprintf(w, string(json))
}

// deleteIngredient is called by IngredientHandler for DELETEs and deletes an ingredient
func deleteIngredient(name string, w http.ResponseWriter) {
	err := models.DeleteIngredient(name)
	if err != nil {
		log.Printf("Error deleting ingredient %v: %v\n", name, err.Error())
		JSONError(w, err, 500)
		return
	}

	json, _ := json.Marshal(Ingredient{Ingredient: name})
	fmt.Fprintf(w, string(json))
}

// IngredientListHandler outputs the list of all the ingredients
func IngredientListHandler(w http.ResponseWriter, r *http.Request) {
	list, err := models.GetAllIngredients()
	if err != nil {
		log.Printf("Error getting ingredient list: %v\n", err.Error())
		JSONError(w, err, 500)
		return
	}

	json, _ := json.Marshal(list)
	fmt.Fprintf(w, string(json))
}
