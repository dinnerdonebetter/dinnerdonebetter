/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $ValidIngredientGroup = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    description: {
      type: 'string',
    },
    id: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    members: {
      type: 'array',
      contains: {
        type: 'ValidIngredientGroupMember',
      },
    },
    name: {
      type: 'string',
    },
    slug: {
      type: 'string',
    },
  },
} as const;
