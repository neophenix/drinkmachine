// Package hw contains the connections to the external hardware
package hw

import (
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/hd44780"
	_ "github.com/kidoman/embd/host/rpi" // load raspberry pi driver
	"log"
	"math"
)

// hd is the driver / connection to the LCD
// this can be nil, as I've found out, if the LCD becomes unplugged, so funcs should check for that.
// I imagine it would be nil if it dies too, so we do want to continue since the web app would still work
var hd *hd44780.HD44780

// CloseDisplay turns off the display and calls embd.CloseI2C() and hd44780.Close()
func CloseDisplay() {
	DisplayToggle(false)
	if hd != nil {
		hd.Close()
	}
	log.Println("Closing I2C")
	embd.CloseI2C()
}

// DisplayToggle turns the display and backlight on / off
func DisplayToggle(on bool) {
	if hd == nil {
		return
	}

	if on {
		hd.BacklightOn()
		hd.DisplayOn()
	} else {
		hd.BacklightOff()
		hd.DisplayOff()
	}
}

// BacklightToggle turns the backlight on / off
func BacklightToggle(on bool) {
	if hd == nil {
		return
	}

	if on {
		hd.BacklightOn()
	} else {
		hd.BacklightOff()
	}
}

// WriteString writes a string to the display on the given line and starting at the given col
// if col == -1 we try to center the string
func WriteString(msg string, line int, col int) {
	if hd == nil {
		return
	}

	if col == -1 {
		col = int(math.Floor(float64((16 - len(msg)) / 2)))
	}

	if len(msg) > 16 {
		msg = msg[0:16]
	}

	hd.SetCursor(col, line)
	for _, b := range []byte(msg) {
		hd.WriteChar(b)
	}
}

// ClearLCD just wraps the HD44780.Clear()
func ClearLCD() {
	hd.Clear()
	hd.SetMode(hd44780.TwoLine)
}

// InitializeLCD inits the I2C interface as well as sets up our LCD driver with the appropriate address
// it then clears the screen, turns it on, and turns the backlight on
func InitializeLCD() {
	if err := embd.InitI2C(); err != nil {
		panic(err)
	}

	bus := embd.NewI2CBus(1)

	hd, _ = hd44780.NewI2C(
		bus,
		0x27,
		hd44780.PCF8574PinMap,
		hd44780.RowAddress16Col,
		hd44780.TwoLine,
	)

	if hd != nil {
		hd.Clear()
		hd.DisplayOn()
		hd.BacklightOn()
	}
}
