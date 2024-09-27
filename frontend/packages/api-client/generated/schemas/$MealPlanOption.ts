/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $MealPlanOption = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    assignedCook: {
      type: 'string',
    },
    assignedDishwasher: {
      type: 'string',
    },
    belongsToMealPlanEvent: {
      type: 'string',
    },
    chosen: {
      type: 'boolean',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    id: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    meal: {
      type: 'Meal',
    },
    mealScale: {
      type: 'number',
      format: 'double',
    },
    notes: {
      type: 'string',
    },
    tieBroken: {
      type: 'boolean',
    },
    votes: {
      type: 'array',
      contains: {
        type: 'MealPlanOptionVote',
      },
    },
  },
} as const;
