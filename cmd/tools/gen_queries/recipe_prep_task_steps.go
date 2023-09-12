package main

const recipePrepTaskStepsTableName = "recipe_prep_task_steps"

var recipePrepTaskStepsColumns = []string{
	"id",
	"satisfies_recipe_step",
	"belongs_to_recipe_step",
	"belongs_to_recipe_prep_task",
}
