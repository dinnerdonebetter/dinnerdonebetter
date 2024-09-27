/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $UserIngredientPreference = {
  properties: {
    allergy: {
      type: 'boolean',
    },
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    belongsToUser: {
      type: 'string',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    id: {
      type: 'string',
    },
    ingredient: {
      type: 'ValidIngredient',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    notes: {
      type: 'string',
    },
    rating: {
      type: 'number',
      format: 'int32',
    },
  },
} as const;
