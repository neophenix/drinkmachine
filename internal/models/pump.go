package models

import (
	"database/sql"
)

// Pump is the model of the db table pumps
// it is similar to hw.Pump except this should match types to the db
type Pump struct {
	ID         int            `sql:"pump_id"`
	FlowRate   float32        `sql:"flow_rate"`
	Ingredient sql.NullString `sql:"ingredient"`
	Pin        string         `sql:"gpio_pin"`
}

// GetPump returns the details of a single pump
func GetPump(id int) (*Pump, error) {
	var p Pump
	err := DB.QueryRow("select pump_id, flow_rate, ingredient, gpio_pin from pumps where pump_id = ?", id).Scan(&p.ID, &p.FlowRate, &p.Ingredient, &p.Pin)
	return &p, err
}

// GetAllPumps returns an ordered list of all the pumps in the DB
func GetAllPumps() ([]*Pump, error) {
	var pumps []*Pump

	rows, err := DB.Query("select pump_id, flow_rate, ingredient, gpio_pin from pumps order by pump_id")
	if err != nil {
		return pumps, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Pump
		err = rows.Scan(&p.ID, &p.FlowRate, &p.Ingredient, &p.Pin)
		if err != nil {
			return pumps, err
		}

		pumps = append(pumps, &p)
	}
	err = rows.Err()
	return pumps, err
}

// UpdatePump allows the user to update either the flow rate or the attached ingredient
func UpdatePump(id int, flowrate float32, ingredient string) error {
	if flowrate == 0 {
		flowrate = 100 // in case someone forgets to set it, pump spec is 100 ml / min
	}

	_, err := DB.Exec("update pumps set flow_rate = ?, ingredient = ? where pump_id = ?", flowrate, ingredient, id)
	return err
}
