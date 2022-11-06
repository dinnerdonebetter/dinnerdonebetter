// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: meal_plan_option_votes_create.sql

package generated

import (
	"context"
)

const CreateMealPlanOptionVote = `-- name: CreateMealPlanOptionVote :exec
INSERT INTO meal_plan_option_votes (id,rank,abstain,notes,by_user,belongs_to_meal_plan_option) VALUES ($1,$2,$3,$4,$5,$6)
`

type CreateMealPlanOptionVoteParams struct {
	ID                      string
	Rank                    int32
	Abstain                 bool
	Notes                   string
	ByUser                  string
	BelongsToMealPlanOption string
}

func (q *Queries) CreateMealPlanOptionVote(ctx context.Context, arg *CreateMealPlanOptionVoteParams) error {
	_, err := q.db.ExecContext(ctx, CreateMealPlanOptionVote,
		arg.ID,
		arg.Rank,
		arg.Abstain,
		arg.Notes,
		arg.ByUser,
		arg.BelongsToMealPlanOption,
	)
	return err
}
