/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeCreationRequestInput = {
  properties: {
    alsoCreateMeal: {
      type: 'boolean',
    },
    description: {
      type: 'string',
    },
    eligibleForMeals: {
      type: 'boolean',
    },
    inspiredByRecipeID: {
      type: 'string',
    },
    maximumEstimatedPortions: {
      type: 'number',
      format: 'double',
    },
    minimumEstimatedPortions: {
      type: 'number',
      format: 'double',
    },
    name: {
      type: 'string',
    },
    pluralPortionName: {
      type: 'string',
    },
    portionName: {
      type: 'string',
    },
    prepTasks: {
      type: 'array',
      contains: {
        type: 'RecipePrepTaskWithinRecipeCreationRequestInput',
      },
    },
    sealOfApproval: {
      type: 'boolean',
    },
    slug: {
      type: 'string',
    },
    source: {
      type: 'string',
    },
    steps: {
      type: 'array',
      contains: {
        type: 'RecipeStepCreationRequestInput',
      },
    },
    yieldsComponentType: {
      type: 'string',
    },
  },
} as const;
