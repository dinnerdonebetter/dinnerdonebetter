package main

import (
	"github.com/cristalhq/builq"
)

const (
	mealPlanOptionVotesTableName = "meal_plan_option_votes"
)

var mealPlanOptionVotesColumns = []string{
	idColumn,
	"rank",
	"abstain",
	notesColumn,
	"by_user",
	"belongs_to_meal_plan_option",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildMealPlanOptionVotesQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}
