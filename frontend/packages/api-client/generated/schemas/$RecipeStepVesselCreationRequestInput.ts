/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepVesselCreationRequestInput = {
  properties: {
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
    productOfRecipeStepIndex: {
      type: 'number',
      format: 'int64',
    },
    productOfRecipeStepProductIndex: {
      type: 'number',
      format: 'int64',
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
