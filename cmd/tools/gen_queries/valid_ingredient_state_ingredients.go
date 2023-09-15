package main

const validIngredientStateIngredientsTableName = "valid_ingredient_state_ingredients"

var validIngredientStateIngredientsColumns = []string{
	idColumn,
	"valid_ingredient",
	"valid_ingredient_state",
	"notes",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientStateIngredientsQueries() []*Query {
	return []*Query{
		//
	}
}
