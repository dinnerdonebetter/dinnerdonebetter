/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepCompletionConditionCreationRequestInput = {
  properties: {
    belongsToRecipeStep: {
      type: 'string',
    },
    ingredientState: {
      type: 'string',
    },
    ingredients: {
      type: 'array',
      contains: {
        type: 'number',
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
