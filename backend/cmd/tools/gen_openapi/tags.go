package main

import (
	"regexp"
)

var routeParamRegex = regexp.MustCompile(`\{[a-zA-Z\d]+\}`)

var tagReplacements = map[string]string{
	"steps":                 "recipe_steps",
	"prep_tasks":            "recipe_prep_tasks",
	"completion_conditions": "recipe_step_completion_conditions",
	"events":                "meal_plan_events",
	"members":               "household_members",
	"ratings":               "recipe_ratings",
	"tasks":                 "meal_plan_tasks",
	"trigger_events":        "webhook_trigger_events",
	"votes":                 "meal_plan_option_votes",
	"vessels":               "recipe_step_vessels",
	"instruments":           "household_instrument_ownerships",
	"configurations":        "service_setting_configurations",
	"settings":              "service_settings",
	"options":               "meal_plan_options",
	"products":              "recipe_step_products",
	"ingredients":           "recipe_step_ingredients",
	"oauth2_clients":        "oauth2",
}
var tagDescriptions = map[string]string{
	"oauth2_clients":                     "",
	"vessels":                            "",
	"user":                               "",
	"valid_vessels":                      "",
	"meal_plan_events":                   "",
	"totp_secret":                        "",
	"valid_preparation_instruments":      "",
	"households":                         "",
	"meal_plan_option_votes":             "",
	"valid_preparation_vessels":          "",
	"webhook_trigger_events":             "",
	"oauth2":                             "",
	"meal_plan_tasks":                    "",
	"service_settings":                   "Operations related to service settings",
	"recipe_step_products":               "",
	"service_setting_configurations":     "Operations related to configuring service settings",
	"valid_ingredients":                  "",
	"email_address":                      "",
	"valid_ingredient_state_ingredients": "",
	"valid_instruments":                  "",
	"auth":                               "",
	"recipe_step_completion_conditions":  "",
	"user_ingredient_preferences":        "",
	"valid_measurement_units":            "",
	"household_invitations":              "",
	"household_instrument_ownerships":    "",
	"household":                          "",
	"valid_ingredient_states":            "",
	"valid_measurement_conversions":      "",
	"webhooks":                           "",
	"recipe_ratings":                     "",
	"password":                           "",
	"valid_ingredient_groups":            "",
	"users":                              "",
	"audit_log_entries":                  "",
	"recipe_step_ingredients":            "",
	"meal_plans":                         "",
	"meal_plan_options":                  "",
	"recipe_steps":                       "",
	"workers":                            "",
	"invitations":                        "",
	"valid_preparations":                 "",
	"meals":                              "",
	"recipe_prep_tasks":                  "",
	"admin":                              "",
	"household_members":                  "",
	"user_notifications":                 "",
	"permissions":                        "",
	"valid_ingredient_preparations":      "",
	"grocery_list_items":                 "",
	"recipes":                            "",
	"valid_ingredient_measurement_units": "",
	"recipe_step_vessels":                "",
}
