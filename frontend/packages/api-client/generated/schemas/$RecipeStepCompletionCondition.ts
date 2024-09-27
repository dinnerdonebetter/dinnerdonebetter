/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepCompletionCondition = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToRecipeStep: {
      type: 'string',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    id: {
      type: 'string',
    },
    ingredientState: {
      type: 'ValidIngredientState',
    },
    ingredients: {
      type: 'array',
      contains: {
        type: 'RecipeStepCompletionConditionIngredient',
      },
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    notes: {
      type: 'string',
    },
    optional: {
      type: 'boolean',
    },
  },
} as const;
