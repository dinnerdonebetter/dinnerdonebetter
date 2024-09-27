/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeStepIngredientCreationRequestInput = {
  properties: {
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
    productOfRecipeStepIndex: {
      type: 'number',
      format: 'int64',
    },
    productOfRecipeStepProductIndex: {
      type: 'number',
      format: 'int64',
    },
    productPercentageToUse: {
      type: 'number',
      format: 'double',
    },
    quantityNotes: {
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
