package main

const validIngredientStatesTableName = "valid_ingredient_states"

var validIngredientStatesColumns = []string{
	idColumn,
	"name",
	"past_tense",
	"slug",
	"description",
	"icon_path",
	"attribute_type",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientStatesQueries() []*Query {
	return []*Query{
		//
	}
}
