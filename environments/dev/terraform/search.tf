# how to use these resources: https://github.com/philippe-vandermoere/terraform-provider-algolia/blob/11d9e162be54c66c92376ae5647f7f3bd675755a/examples/main.tf

locals {
  default_algolia_ranking_criteria = ["typo", "words", "filters", "proximity", "attribute", "exact", "custom"]

  # common rankings
  name_and_description_custom_ranking = [
    "desc(name)",
    "desc(description)",
  ]

  name_plural_name_and_description_custom_ranking = [
    "desc(name)",
    "desc(pluralName)",
    "desc(description)",
  ]

  name_past_tense_and_description_custom_ranking = [
    "desc(name)",
    "desc(pastTense)",
    "desc(description)",
  ]
}

resource "algolia_index" "recipes_index" {
  name = "recipes"

  ranking = concat(local.name_and_description_custom_ranking, local.default_algolia_ranking_criteria)

  searchable_attributes = [
    "name",
    "description",
  ]
}

resource "algolia_index" "meals_index" {
  name = "meals"

  ranking = concat(local.name_and_description_custom_ranking, local.default_algolia_ranking_criteria)

  searchable_attributes = [
    "name",
    "description",
  ]
}

resource "algolia_index" "valid_ingredients_index" {
  name = "valid_ingredients"

  ranking = concat(local.name_plural_name_and_description_custom_ranking, local.default_algolia_ranking_criteria)

  searchable_attributes = [
    "name",
    "pluralName",
    "description",
  ]
}

resource "algolia_index" "valid_instruments_index" {
  name = "valid_instruments"

  ranking = concat(local.name_plural_name_and_description_custom_ranking, local.default_algolia_ranking_criteria)

  searchable_attributes = [
    "name",
    "pluralName",
    "description",
  ]
}

resource "algolia_index" "valid_measurement_units_index" {
  name = "valid_measurement_units"

  ranking = concat(local.name_plural_name_and_description_custom_ranking, local.default_algolia_ranking_criteria)

  searchable_attributes = [
    "name",
    "pluralName",
    "description",
  ]
}

resource "algolia_index" "valid_preparations_index" {
  name = "valid_preparations"

  ranking = concat(local.name_past_tense_and_description_custom_ranking, local.default_algolia_ranking_criteria)

  searchable_attributes = [
    "name",
    "pastTense",
    "description",
  ]
}

resource "algolia_index" "valid_ingredient_states_index" {
  name = "valid_ingredient_states"

  ranking = concat(local.name_past_tense_and_description_custom_ranking, local.default_algolia_ranking_criteria)

  searchable_attributes = [
    "name",
    "pastTense",
    "description",
  ]
}
