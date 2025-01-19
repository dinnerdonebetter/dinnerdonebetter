# how to use these resources: https://github.com/philippe-vandermoere/terraform-provider-algolia/blob/11d9e162be54c66c92376ae5647f7f3bd675755a/examples/main.tf

resource "algolia_index" "recipes_index" {
  name = "recipes"
}

resource "algolia_index" "meals_index" {
  name = "meals"
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

resource "algolia_index" "valid_ingredient_states_index" {
  name = "valid_ingredient_states"
}

resource "algolia_index" "users_index" {
  name = "users"
}
