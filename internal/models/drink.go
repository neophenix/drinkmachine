package models

import (
	"errors"
)

// Drink is the model of the drinks table
type Drink struct {
	ID          int64              `sql:"drink_id" json:"id"`
	Name        string             `sql:"name" json:"name"`
	Notes       string             `sql:"notes" json:"notes"`
	Ingredients []*DrinkIngredient `json:"ingredients"`
}

// DrinkIngredient is the model of the drink_ingredients table
type DrinkIngredient struct {
	Ingredient string  `sql:"ingredient" json:"ingredient"`
	Amount     float32 `sql:"amount" json:"amount"`
	Units      string  `sql:"units" json:"units"`
	Dispense   bool    `sql:"dispense" json:"dispense"`
	DrinkID    int64   `sql:"drink_id" json:"drink_id"`
}

// GetDrink returns the drink + ingredients to the user, it looks it up by id or name in that order
func GetDrink(id int64, name string) (*Drink, error) {
	var d Drink
	var err error

	// This is probably pointless as we will always look up by id, but I wrote it and not removing it in the off
	// chance I find out I need it.  Maybe in 2032 I'll remove it
	if id != 0 {
		err = DB.QueryRow("select drink_id, name, notes from drinks where drink_id = ?", id).Scan(&d.ID, &d.Name, &d.Notes)
	} else if name != "" {
		err = DB.QueryRow("select drink_id, name, notes from drinks where lower(name) = lower(?)", name).Scan(&d.ID, &d.Name, &d.Notes)
	} else {
		err = errors.New("Neither ID or Name was passed to lookup drink")
	}

	if err != nil {
		return nil, err
	}

	rows, err := DB.Query("select ingredient, amount, units, dispense, drink_id from drink_ingredients where drink_id = ?", d.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var di DrinkIngredient
		err = rows.Scan(&di.Ingredient, &di.Amount, &di.Units, &di.Dispense, &di.DrinkID)
		if err != nil {
			return nil, err
		}

		d.Ingredients = append(d.Ingredients, &di)
	}

	return &d, nil
}

// CreateDrink puts a drink in the drinks table, ingredients are added elsewhere
func CreateDrink(name string, notes string) (*Drink, error) {
	res, err := DB.Exec("insert into drinks (name, notes) values (?,?)", name, notes)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Drink{ID: id, Name: name, Notes: notes}, nil
}

// ClearIngredients deletes all the drink_ingredients rows for a drink, this would be in prep to add new ones, that is far
// easier than any complicated comparison of what we want vs what we have
func (d *Drink) ClearIngredients() error {
	_, err := DB.Exec("delete from drink_ingredients where drink_id = ?", d.ID)
	return err
}

// AddIngredient adds an entry to the drink_ingredients table
func (d *Drink) AddIngredient(ingredient string, amount float32, units string, dispense bool) error {
	_, err := DB.Exec("insert into drink_ingredients (ingredient, amount, units, dispense, drink_id) values (?,?,?,?,?)", ingredient, amount, units, dispense, d.ID)
	if err != nil {
		return err
	}

	d.Ingredients = append(d.Ingredients, &DrinkIngredient{Ingredient: ingredient, Amount: amount, Units: units, Dispense: dispense, DrinkID: d.ID})
	return nil
}

// UpdateDrink lets the user update the name / notes of a drink, for ingredient updates call ClearIngredients / AddIngredient
func UpdateDrink(id int64, name string, notes string) (*Drink, error) {
	_, err := DB.Exec("update drinks set name = ?, notes = ? where drink_id = ?", name, notes, id)
	if err != nil {
		return nil, err
	}

	return &Drink{ID: id, Name: name, Notes: notes}, nil
}

// DeleteDrink removes the drink and its ingredients from the db
func DeleteDrink(id int64) error {
	_, err := DB.Exec("delete from drink_ingredients where drink_id = ?", id)
	if err != nil {
		return err
	}

	_, err = DB.Exec("delete from drinks where drink_id = ?", id)
	return err
}

// GetAllDrinks returns all the drinks in the system ordered by name
// if simple is true, don't populate the ingredients
func GetAllDrinks(simple bool) ([]*Drink, error) {
	var drinks []*Drink

	rows, err := DB.Query("select drink_id, name, notes from drinks order by lower(name)")
	if err != nil {
		return drinks, err
	}
	defer rows.Close()
	for rows.Next() {
		var d Drink
		err = rows.Scan(&d.ID, &d.Name, &d.Notes)
		if err != nil {
			return drinks, err
		}

		if simple != true {
			irows, err := DB.Query("select ingredient, amount, units, dispense, drink_id from drink_ingredients where drink_id = ?", d.ID)
			if err != nil {
				return drinks, err
			}
			defer irows.Close()
			for irows.Next() {
				var di DrinkIngredient
				err = irows.Scan(&di.Ingredient, &di.Amount, &di.Units, &di.Dispense, &di.DrinkID)
				if err != nil {
					return drinks, err
				}

				d.Ingredients = append(d.Ingredients, &di)
			}

		}

		drinks = append(drinks, &d)
	}

	return drinks, nil
}
