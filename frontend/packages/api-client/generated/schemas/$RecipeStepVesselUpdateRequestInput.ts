/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepVesselUpdateRequestInput = {
  properties: {
    belongsToRecipeStep: {
      type: 'string',
    },
    maximumQuantity: {
      type: 'number',
      format: 'int64',
    },
    minimumQuantity: {
      type: 'number',
      format: 'int64',
    },
    name: {
      type: 'string',
    },
    notes: {
      type: 'string',
    },
    recipeStepProductID: {
      type: 'string',
    },
    unavailableAfterStep: {
      type: 'boolean',
    },
    vesselID: {
      type: 'string',
    },
    vesselPreposition: {
      type: 'string',
    },
  },
} as const;
