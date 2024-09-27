/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $MealPlanTask = {
  properties: {
    assignedToUser: {
      type: 'string',
    },
    completedAt: {
      type: 'string',
      format: 'date-time',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    creationExplanation: {
      type: 'string',
    },
    id: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    mealPlanOption: {
      type: 'MealPlanOption',
    },
    recipePrepTask: {
      type: 'RecipePrepTask',
    },
    status: {
      type: 'string',
    },
    statusExplanation: {
      type: 'string',
    },
  },
} as const;
