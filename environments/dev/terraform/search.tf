# how to use these resources: https://github.com/philippe-vandermoere/terraform-provider-algolia/blob/11d9e162be54c66c92376ae5647f7f3bd675755a/examples/main.tf

locals {
  default_algolia_ranking_criteria = ["custom", "exact", "typo", "words", "attribute", "filters", "proximity"]
}

resource "algolia_index" "recipes_index" {
  name = "recipes"

  attributes_config {
    searchable_attributes = [
      "name",
      "description",
    ]
  }
}

resource "algolia_index" "meals_index" {
  name = "meals"

  attributes_config {
    searchable_attributes = [
      "name",
      "description",
    ]
  }
}

resource "algolia_index" "valid_ingredients_index" {
  name = "valid_ingredients"

  attributes_config {
    searchable_attributes = [
      "name",
      "pluralName",
      "description",
    ]
  }
}

resource "algolia_index" "valid_instruments_index" {
  name = "valid_instruments"

  attributes_config {
    searchable_attributes = [
      "name",
      "pluralName",
      "description",
    ]
  }
}

resource "algolia_index" "valid_measurement_units_index" {
  name = "valid_measurement_units"

  attributes_config {
    searchable_attributes = [
      "name",
      "pluralName",
      "description",
    ]
  }
}

resource "algolia_index" "valid_preparations_index" {
  name = "valid_preparations"

  attributes_config {
    searchable_attributes = [
      "name",
      "pastTense",
      "description",
    ]
  }
}

resource "algolia_index" "valid_ingredient_states_index" {
  name = "valid_ingredient_states"

  attributes_config {
    searchable_attributes = [
      "name",
      "pastTense",
      "description",
    ]
  }
}

resource "algolia_index" "users_index" {
  name = "users"

  attributes_config {
    searchable_attributes = [
      "username",
    ]
  }
}
