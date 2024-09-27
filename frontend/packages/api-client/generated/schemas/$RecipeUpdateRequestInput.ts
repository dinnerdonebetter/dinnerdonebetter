/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $RecipeUpdateRequestInput = {
  properties: {
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
    sealOfApproval: {
      type: 'boolean',
    },
    slug: {
      type: 'string',
    },
    source: {
      type: 'string',
    },
    yieldsComponentType: {
      type: 'string',
    },
  },
} as const;
