package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/neophenix/drinkmachine/internal/hw"
	"github.com/neophenix/drinkmachine/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Pump is the JSON version of a pump, every package gets a pump
type Pump struct {
	ID         int     `json:"id"`
	FlowRate   float32 `json:"flow_rate"`
	Ingredient string  `json:"ingredient"`
}

// PumpHandler handles:
//  GET - output info about a pump
//  PUT - update flow rate and / or ingredient
func PumpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// path should be something like /api/pump/X
	parts := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(parts[3])

	// return the pump info
	if r.Method == "GET" {
		pump, err := models.GetPump(id)
		if err != nil {
			log.Printf("Error getting pump: %v\n", err.Error())
			JSONError(w, err, 500)
			return
		}
		json, _ := json.Marshal(Pump{ID: pump.ID, FlowRate: pump.FlowRate, Ingredient: pump.Ingredient.String})
		fmt.Fprintf(w, string(json))
	}

	// update and return the pump info
	if r.Method == "PUT" {
		buf := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
		buf.ReadFrom(r.Body)

		var put Pump
		err := json.Unmarshal(buf.Bytes(), &put)
		if err != nil {
			log.Printf("Error reading incoming data: %v\n", err.Error())
			JSONError(w, err, 400)
			return
		}

		err = models.UpdatePump(id, put.FlowRate, put.Ingredient)
		if err != nil {
			log.Printf("Error updating pump %v: %v\n", id, err.Error())
			JSONError(w, err, 500)
			return
		}

		pump, err := models.GetPump(id)
		if err != nil {
			log.Printf("Error getting pump: %v\n", err.Error())
			JSONError(w, err, 500)
			return
		}

		// update the hw package details for the given pump so we don't have to re-init or anything
		hw.UpdatePump(id, pump.FlowRate, pump.Ingredient.String)

		json, _ := json.Marshal(put)
		fmt.Fprintf(w, string(json))
	}
}

// PumpListHandler outputs the list of all pumps and their info
// This actually pulls the list from the hw package, not the DB as the hw would be what the system is
// basing its decisions off of and seems more useful
func PumpListHandler(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(hw.Pumps)
	fmt.Fprintf(w, string(json))
}
