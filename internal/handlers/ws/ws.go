// Package ws is for our websocket handler
package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/neophenix/drinkmachine/internal/hw"
	"github.com/neophenix/drinkmachine/internal/models"
	"log"
	"net/http"
	"strings"
	"time"
)

// IncomingMessage is the incoming message from the UI that tells us what drink to make
type IncomingMessage struct {
	Action  string                 `json:"action"`  // what we are doing, which should always be "make_drink" at least for now
	ID      int                    `json:"id"`      // the id of the drink we are going to make, or pump to use
	Options IncomingMessageOptions `json:"options"` // options of the action
}

// IncomingMessageOptions hold options for actions from the frontend
//  * Duration - how long to run a pump when run_pump is passed
type IncomingMessageOptions struct {
	Duration int  `json:"duration"` // how long to run the pump
	Override bool `json:"override"` // override flag to pour a drink with missing ingredients, maybe something else in the future
}

// OutgoingMessage hold status updates we want to communicate to the user
type OutgoingMessage struct {
	ID      string `json:"id"`      // an id for the UI to use to identify messages coming to it
	Type    string `json:"type"`    // for when we actually start pouring, the type of message
	Message string `json:"message"` // text we want displayed to the user
	Success bool   `json:"success"` // tell the frontend whether this is a good or bad message
}

// stopRunning is set to true when a stop command is received from the frontend, maybe
// some ingredient ran out or the user picked the wrong drink.
var stopRunning = false

// default upgrader
var upgrader = websocket.Upgrader{}

// Handler is the websocket handler, the only thing we do over the socket is make the
// drink since we want to report back times, completed ingredients, etc to the user
func Handler(w http.ResponseWriter, r *http.Request) {
	// upgrade to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v\n", err.Error())
		w.WriteHeader(500)
		return
	}
	defer conn.Close()

	for {
		// read the raw (encoded) message from the connection
		msgtype, encodedmsg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading ws message: %v\n", err.Error())
			jsonb, _ := json.Marshal(OutgoingMessage{Type: "error", Message: "Error reading request " + err.Error(), Success: false})
			conn.WriteMessage(msgtype, jsonb)
			break
		}

		// Log what we got for later debugging and unmarshal it to our struct
		var msg IncomingMessage
		err = json.Unmarshal(encodedmsg, &msg)
		log.Printf("ws recv: %v\n", msg)
		if err != nil {
			log.Printf("error decoding ws message: %v\n", err.Error())
			jsonb, _ := json.Marshal(OutgoingMessage{Type: "error", Message: "Did not understand request " + err.Error(), Success: false})
			conn.WriteMessage(msgtype, jsonb)
			break
		}

		switch msg.Action {
		case "make_drink":
			locked := hw.LockPumps()
			if !locked {
				// don't want to queue up users with a normal lock, so tell them to come back later
				jsonb, _ := json.Marshal(OutgoingMessage{Type: "error", Message: "System currently in use, please wait and try again later.", Success: false})
				conn.WriteMessage(msgtype, jsonb)
				break
			}
			// make sure our stop flag is false, user can hit stop from here on out
			stopRunning = false
			makeDrink(conn, msgtype, msg)
			// free up the lock for the next person
			hw.UnlockPumps()
		case "run_pump":
			locked := hw.LockPumps()
			if !locked {
				// don't want to queue up users with a normal lock, so tell them to come back later
				jsonb, _ := json.Marshal(OutgoingMessage{Type: "error", Message: "System currently in use, please wait and try again later.", Success: false})
				conn.WriteMessage(msgtype, jsonb)
				break
			}
			runPump(conn, msgtype, msg)
			// free up the lock for the next person
			hw.UnlockPumps()
		case "stop":
			// just set the stop flag to true, the current running loop should stop within a second
			stopRunning = true
		default:
			jsonb, _ := json.Marshal(OutgoingMessage{Type: "error", Message: "Request not understood", Success: false})
			conn.WriteMessage(msgtype, jsonb)
		}
	}
}

// drinkPump holds the info on a pump we are going to use for the drink
type drinkPump struct {
	Pump       *hw.Pump                // the actual connection to the pump
	Ingredient *models.DrinkIngredient // all the info about the ingredient will be in here
	Run        bool                    // Will initially be true, then set to false when we don't need the pump anymore
}

// makeDrink handles figuring out what pumps to run and kicking them off, communicating back to the user and updating LCD
func makeDrink(conn *websocket.Conn, msgtype int, msg IncomingMessage) {
	drink, err := models.GetDrink(int64(msg.ID), "")
	if err != nil {
		jsonb, _ := json.Marshal(OutgoingMessage{Type: "error", Message: "Could not find that drink", Success: false})
		conn.WriteMessage(msgtype, jsonb)
		return
	}

	// make a lookup that collects all "like pumps" so we can find them easier later
	var lookup map[string][]*hw.Pump = make(map[string][]*hw.Pump)
	for _, p := range hw.Pumps {
		lookup[strings.ToLower(p.Ingredient)] = append(lookup[strings.ToLower(p.Ingredient)], p)
	}

	var missing []string        // any missing ingredients
	var dispensing []string     // ingredients we will actually dispense
	var pumpsToUse []*drinkPump // list of all the pumps we are going to need
	for _, i := range drink.Ingredients {
		// only look at things we can pump
		if i.Dispense {
			list, ok := lookup[strings.ToLower(i.Ingredient)]
			if ok {
				for _, p := range list {
					pumpsToUse = append(pumpsToUse, &drinkPump{Pump: p, Ingredient: i, Run: true})
				}
				dispensing = append(dispensing, fmt.Sprintf("%v %v %v", i.Amount, i.Units, i.Ingredient))
			} else {
				// we don't have this ingredient so add it here and we will tell the user in a moment
				missing = append(missing, fmt.Sprintf("%v %v %v", i.Amount, i.Units, i.Ingredient))
			}
		}
	}

	// tell the user that they wanted a drink that we don't have the right hookups for
	if len(missing) > 0 {
		// If the user has chosen to override the missing indredients, pass them back up as a finish item,
		// so they have a full list of things to add when we are done
		if msg.Options.Override {
			for _, ingredient := range missing {
				jsonb, _ := json.Marshal(OutgoingMessage{Type: "finish", Message: ingredient, Success: true})
				conn.WriteMessage(msgtype, jsonb)
			}
		} else {
			jsonb, _ := json.Marshal(OutgoingMessage{Type: "missing", Message: strings.Join(missing, "<br/>"), Success: false})
			conn.WriteMessage(msgtype, jsonb)
			return
		}
	}

	// lets us bail and cleanup in the event of an error, stop all the pumps
	defer hw.StopAllPumps()

	// give the UI a list of ingredients we are going to pour
	for _, ingredient := range dispensing {
		jsonb, _ := json.Marshal(OutgoingMessage{Type: "pouring", Message: ingredient, Success: true})
		conn.WriteMessage(msgtype, jsonb)
	}

	// now send up anything we aren't going to dispense, it will hide this until we are done
	for _, i := range drink.Ingredients {
		if !i.Dispense {
			message := fmt.Sprintf("%v %v %v", i.Amount, i.Units, i.Ingredient)
			jsonb, _ := json.Marshal(OutgoingMessage{ID: i.Ingredient, Type: "finish", Message: message, Success: true})
			conn.WriteMessage(msgtype, jsonb)
		}
	}

	// and lastly any notes on the drink
	if drink.Notes != "" {
		jsonb, _ := json.Marshal(OutgoingMessage{Type: "notes", Message: drink.Notes, Success: true})
		conn.WriteMessage(msgtype, jsonb)
	}

	// display something on the LCD since I have that hooked up
	hw.DisplayToggle(true)
	defer hw.DisplayToggle(false)
	hw.ClearLCD()

	hw.WriteString("Pouring in...", 0, -1)
	countdown := 3
	for countdown > 0 {
		hw.WriteString(fmt.Sprintf("%v", countdown), 1, -1)
		time.Sleep(1 * time.Second)
		countdown--
	}

	// We made it through any error checking, so now that we are going to make the drink send a message to the UI to show the appropriate info
	jsonb, _ := json.Marshal(OutgoingMessage{Type: "starting", Message: "", Success: true})
	conn.WriteMessage(msgtype, jsonb)

	hw.ClearLCD()
	hw.WriteString("Pouring", 0, -1)
	hw.WriteString(drink.Name, 1, -1)

	var timeRemaining float64 = 1
	for timeRemaining > 0 && !stopRunning {
		// we set these each loop so reset them here to a known state
		timeRemaining = 0

		for _, p := range pumpsToUse {
			if p.Run {
				runTime, err := p.Pump.Run(p.Ingredient.Amount, p.Ingredient.Units, p.Ingredient.Ingredient)
				if err != nil {
					// this will hopefully only (never?) happen when we first try to run
					jsonb, _ := json.Marshal(OutgoingMessage{Type: "error", Message: fmt.Sprintf("Error running pump %v %v", p.Pump.ID, err.Error()), Success: false})
					conn.WriteMessage(msgtype, jsonb)
					return
				}

				// don't need to run anymore, so this pump has finished
				if runTime <= 0 {
					p.Run = false
					// tell the UI we are done with this one, multiple of the same will be ok
					jsonb, _ := json.Marshal(OutgoingMessage{ID: p.Ingredient.Ingredient, Type: "pour_complete", Success: true})
					conn.WriteMessage(msgtype, jsonb)
				} else {
					// use the greater of runTime or timeRemaining so we keep looping and can report back to the user
					if runTime > timeRemaining {
						timeRemaining = runTime
					}
				}
			}
		}

		// report back how much time we have left
		msg := fmt.Sprintf("%v", timeRemaining)
		jsonb, _ := json.Marshal(OutgoingMessage{Type: "time_remaining", Message: msg, Success: true})
		conn.WriteMessage(msgtype, jsonb)

		if timeRemaining > 0 {
			time.Sleep(1 * time.Second)
		}
	}

	// If we stopped because the user hit stop, shut off all the pumps, and just return back a time_remaining of 0 for now so the UI thinks we finished.
	// Later I'd like to get fancier here and perhaps hand back a list of amounts still missing
	if stopRunning {
		hw.StopAllPumps()
		jsonb, _ := json.Marshal(OutgoingMessage{Type: "time_remaining", Message: "0", Success: true})
		conn.WriteMessage(msgtype, jsonb)
	}
}

func runPump(conn *websocket.Conn, msgtype int, msg IncomingMessage) {
	// 5 seconds seems like a good default duration to run a pump
	duration := 5
	if msg.Options.Duration != 0 {
		duration = msg.Options.Duration
	}

	if msg.ID < 1 || msg.ID > 8 {
		jsonb, _ := json.Marshal(OutgoingMessage{Type: "error", Message: fmt.Sprintf("Invalid Pump number %v", msg.ID), Success: false})
		conn.WriteMessage(msgtype, jsonb)
		return
	}

	pump := hw.Pumps[msg.ID-1]
	pump.RunSeconds(duration)
	jsonb, _ := json.Marshal(OutgoingMessage{Type: "done", Message: "Done", Success: true})
	conn.WriteMessage(msgtype, jsonb)
}
