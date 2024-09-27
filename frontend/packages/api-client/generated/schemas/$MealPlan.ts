/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $MealPlan = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToHousehold: {
      type: 'string',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    createdBy: {
      type: 'string',
    },
    electionMethod: {
      type: 'string',
    },
    events: {
      type: 'array',
      contains: {
        type: 'MealPlanEvent',
      },
    },
    groceryListInitialized: {
      type: 'boolean',
    },
    id: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    notes: {
      type: 'string',
    },
    status: {
      type: 'string',
    },
    tasksCreated: {
      type: 'boolean',
    },
    votingDeadline: {
      type: 'string',
      format: 'date-time',
    },
  },
} as const;
