// Package hw contains the connections to the external hardware
package hw

import (
	"errors"
	"fmt"
	"github.com/kidoman/embd"
	"github.com/neophenix/drinkmachine/internal/models"
	"github.com/neophenix/drinkmachine/internal/utils"
	"log"
	"strings"
	"sync/atomic"
	"time"
)

// lock is a simple counter to use with atomic functions to see if we are the only ones
// trying to access the pumps.  I really wanted TryLock so as to not block and return immediately,
// but without that just needed something very simple to use
var lock int32

// specFlowRate is the spec of the pumps in ml / min
var specFlowRate = float32(100)

// Pump stores all the info we need to know about a pump
type Pump struct {
	ID         int             `json:"id"`         // The pump id from the user's perspective
	FlowRate   float32         `json:"flow_rate"`  // Observed flow rate of the pump, run it for 10s, measure what you got and multiply by 6
	Ingredient string          `json:"ingredient"` // Name of the ingredient hooked up to the pump
	Pin        embd.DigitalPin `json:"-"`          // The GPIO pin object this pump (relay) is connected to
	RunTime    time.Time       `json:"-"`          // The timestamp of when this pump started running
}

// Pumps is our list of pumps in order from ID 1 - X
var Pumps []*Pump

// Run will start or continue a pump running until it has determined it has pumped the
// correct amount of the given ingredient.
//
// returns the seconds the pump has yet to run, and / or error
func (p *Pump) Run(amount float32, units string, ingredient string) (float64, error) {
	if strings.ToLower(ingredient) != strings.ToLower(p.Ingredient) {
		return 0, errors.New("Pump is not connected to " + ingredient)
	}

	amount, err := utils.ConvertToML(amount, units)
	if err != nil {
		return 0, err
	}

	totalRunTime := float64((amount / (p.FlowRate * specFlowRate)) * 60)
	now := time.Now()

	// run the pump, I guess this could be < 1 second, but lets overpour rather than under?
	if totalRunTime > 0 {
		if p.RunTime == (time.Time{}) {
			// we haven't started pumping yet
			p.RunTime = now
			p.Pin.Write(embd.Low)
			return totalRunTime, nil
		}

		// we are pumping, so figure out how much time we have left
		runningTime := now.Sub(p.RunTime).Seconds()
		if runningTime >= totalRunTime {
			// we have run for the correct amount of time, shut it down
			p.Pin.Write(embd.High)
			p.RunTime = time.Time{}
			return 0, nil
		}

		return totalRunTime - runningTime, nil
	}

	// we should never get here but just in case
	return 0, nil
}

// RunSeconds runs the specified pump for X seconds, useful for priming and cleaning.
// this will block with a sleep unlike Run
func (p *Pump) RunSeconds(seconds int) {
	DisplayToggle(true)
	ClearLCD()
	WriteString(fmt.Sprintf("Running Pump %v", p.ID), 0, -1)
	p.Pin.Write(embd.Low)
	for seconds > 0 {
		WriteString(fmt.Sprintf("%v", seconds), 1, -1)
		time.Sleep(1 * time.Second)
		seconds--
	}
	p.Pin.Write(embd.High)
	DisplayToggle(false)
}

// StopAllPumps sends a High signal to all the pumps, stopping them, useful in the event of an error
func StopAllPumps() {
	log.Println("Stopping all pumps")
	for _, p := range Pumps {
		p.Pin.Write(embd.High)
		// reset any internal timers for the pump
		p.RunTime = time.Time{}
	}
}

// UpdatePump called by the API to update the ingredient / flow rate of a pump
func UpdatePump(id int, flowrate float32, ingredient string) {
	// Its easier for a user to just enter the whole number they saw, but we use this as an adjustment so convert to percentage
	if flowrate >= 1 {
		flowrate /= 100
	}

	Pumps[id-1].FlowRate = flowrate
	Pumps[id-1].Ingredient = ingredient
}

// LockPumps tries to set the lock to 1, returning the result of a CompareAndSwap, if it returns true
// then we have the lock
func LockPumps() bool {
	return atomic.CompareAndSwapInt32(&lock, 0, 1)
}

// UnlockPumps tries to set the lock to 0, returning the result of a CompareAndSwap
func UnlockPumps() bool {
	return atomic.CompareAndSwapInt32(&lock, 1, 0)
}

// ClosePumps wraps embd.CloseGPIO for the main function to call when it ends
func ClosePumps() {
	StopAllPumps() // just in case?
	log.Println("Closing GPIO")
	embd.CloseGPIO()
}

// InitializePumps pulls the pump info from the DB and sets up the Pumps list
// it also sets up the GPIO connection to the relay / pump
func InitializePumps() {
	// This gets closed down in ClosePumps() which is called in main
	embd.InitGPIO()

	dbPumps, err := models.GetAllPumps()
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range dbPumps {
		gpio, err := embd.NewDigitalPin(p.Pin)
		if err != nil {
			log.Fatal(err)
		}
		gpio.SetDirection(embd.Out)
		gpio.Write(embd.High)

		if p.FlowRate >= 1 {
			p.FlowRate /= 100
		}

		pump := &Pump{
			ID:         p.ID,
			FlowRate:   p.FlowRate,
			Ingredient: p.Ingredient.String,
			Pin:        gpio,
			RunTime:    time.Time{},
		}

		Pumps = append(Pumps, pump)
	}
}
