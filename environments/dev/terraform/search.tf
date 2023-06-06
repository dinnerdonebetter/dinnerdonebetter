resource "algolia_index" "recipes_index" {
  name = "recipes"
}

resource "algolia_index" "valid_ingredients_index" {
  name = "valid_ingredients"
}

resource "algolia_index" "valid_instruments_index" {
  name = "valid_instruments"
}

resource "algolia_index" "valid_measurement_units_index" {
  name = "valid_measurement_units"
}

resource "algolia_index" "valid_preparations_index" {
  name = "valid_preparations"
}

resource "algolia_index" "valid_ingredient_measurement_units_index" {
  name = "valid_ingredient_measurement_units"
}

