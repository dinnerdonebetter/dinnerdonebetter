/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepCompletionConditionForExistingRecipeCreationRequestInput = {
  properties: {
    belongsToRecipeStep: {
      type: 'string',
    },
    ingredientStateID: {
      type: 'string',
    },
    ingredients: {
      type: 'array',
      contains: {
        type: 'RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput',
      },
    },
    notes: {
      type: 'string',
    },
    optional: {
      type: 'boolean',
    },
  },
} as const;
