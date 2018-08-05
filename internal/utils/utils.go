// Package utils for various utility functions that might be useful to share?
package utils

import (
	"errors"
)

// ConvertToML will take typical drink recipe units and convert them to
// ml as everything we do with the pumps we assume ml units in order
// to calculate how long they need to be on
func ConvertToML(amount float32, unit string) (float32, error) {
	switch unit {
	case "ml":
		return amount, nil
	case "cl":
		return amount * 10, nil
	case "oz":
		return amount * 29.5735, nil
	}

	return 0, errors.New("Unit type " + unit + " not understood")
}
