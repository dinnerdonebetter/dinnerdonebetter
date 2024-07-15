// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: recipe_prep_task_steps.sql

package generated

import (
	"context"
)

const createRecipePrepTaskStep = `-- name: CreateRecipePrepTaskStep :exec

INSERT INTO recipe_prep_task_steps (
	id,
	belongs_to_recipe_step,
	belongs_to_recipe_prep_task,
	satisfies_recipe_step
) VALUES (
	$1,
	$2,
	$3,
	$4
)
`

type CreateRecipePrepTaskStepParams struct {
	ID                      string
	BelongsToRecipeStep     string
	BelongsToRecipePrepTask string
	SatisfiesRecipeStep     bool
}

func (q *Queries) CreateRecipePrepTaskStep(ctx context.Context, db DBTX, arg *CreateRecipePrepTaskStepParams) error {
	_, err := db.ExecContext(ctx, createRecipePrepTaskStep,
		arg.ID,
		arg.BelongsToRecipeStep,
		arg.BelongsToRecipePrepTask,
		arg.SatisfiesRecipeStep,
	)
	return err
}
