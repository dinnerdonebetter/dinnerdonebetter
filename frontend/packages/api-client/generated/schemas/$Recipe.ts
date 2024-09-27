/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $Recipe = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    createdByUser: {
      type: 'string',
    },
    description: {
      type: 'string',
    },
    eligibleForMeals: {
      type: 'boolean',
    },
    id: {
      type: 'string',
    },
    inspiredByRecipeID: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    maximumEstimatedPortions: {
      type: 'number',
      format: 'double',
    },
    media: {
      type: 'array',
      contains: {
        type: 'RecipeMedia',
      },
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
        type: 'RecipePrepTask',
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
        type: 'RecipeStep',
      },
    },
    supportingRecipes: {
      type: 'array',
      contains: {
        type: 'Recipe',
      },
    },
    yieldsComponentType: {
      type: 'string',
    },
  },
} as const;
