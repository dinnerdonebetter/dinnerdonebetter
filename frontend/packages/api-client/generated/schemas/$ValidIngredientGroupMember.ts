/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $ValidIngredientGroupMember = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToGroup: {
      type: 'string',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    id: {
      type: 'string',
    },
    validIngredient: {
      type: 'ValidIngredient',
    },
  },
} as const;
