package models

// Ingredient is the model of the db table ingredients
type Ingredient struct {
	Ingredient string  `sql:"ingredient"`
	Viscosity  float32 `sql:"viscosity"`
}

// GetIngredient returns the ingredient row for the given name
func GetIngredient(name string) (*Ingredient, error) {
	var i Ingredient
	err := DB.QueryRow("select ingredient, viscosity from ingredients where lower(ingredient) = lower(?)", name).Scan(&i.Ingredient, &i.Viscosity)
	return &i, err
}

// GetAllIngredients returns all ingredients sorted by name from the db
func GetAllIngredients() ([]*Ingredient, error) {
	var list []*Ingredient

	rows, err := DB.Query("select ingredient, viscosity from ingredients order by lower(ingredient)")
	if err != nil {
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var i Ingredient
		err = rows.Scan(&i.Ingredient, &i.Viscosity)
		if err != nil {
			return list, err
		}

		list = append(list, &i)
	}
	err = rows.Err()
	return list, err
}

// CreateIngredient puts a new ingredient in the db
func CreateIngredient(ingredient string, viscosity float32) error {
	if viscosity == 0 {
		viscosity = 1 // in case someone forgets to set it
	}

	_, err := DB.Exec("insert into ingredients (ingredient, viscosity) values (?,?)", ingredient, viscosity)
	return err
}

// UpdateIngredient lets a user update the viscosity of the fluid, if we ever need that
func UpdateIngredient(ingredient string, viscosity float32) error {
	if viscosity == 0 {
		viscosity = 1 // in case someone forgets to set it
	}

	_, err := DB.Exec("update ingredients set viscosity = ? where lower(ingredient) = lower(?)", viscosity, ingredient)
	return err
}

// DeleteIngredient lets users cover up their mistakes because they probably typoed something
func DeleteIngredient(ingredient string) error {
	_, err := DB.Exec("delete from ingredients where lower(ingredient) = lower(?)", ingredient)
	return err
}
