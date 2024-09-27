/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepIngredientUpdateRequestInput = {
  properties: {
    belongsToRecipeStep: {
      type: 'string',
    },
    ingredientID: {
      type: 'string',
    },
    ingredientNotes: {
      type: 'string',
    },
    maximumQuantity: {
      type: 'number',
      format: 'double',
    },
    measurementUnitID: {
      type: 'string',
    },
    minimumQuantity: {
      type: 'number',
      format: 'double',
    },
    name: {
      type: 'string',
    },
    optionIndex: {
      type: 'number',
      format: 'int64',
    },
    optional: {
      type: 'boolean',
    },
    productOfRecipeID: {
      type: 'string',
    },
    productPercentageToUse: {
      type: 'number',
      format: 'double',
    },
    quantityNotes: {
      type: 'string',
    },
    recipeStepProductID: {
      type: 'string',
    },
    toTaste: {
      type: 'boolean',
    },
    vesselIndex: {
      type: 'number',
      format: 'int64',
    },
  },
} as const;
