/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $ValidIngredientMeasurementUnit = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
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
    maximumAllowableQuantity: {
      type: 'number',
      format: 'double',
    },
    measurementUnit: {
      type: 'ValidMeasurementUnit',
    },
    minimumAllowableQuantity: {
      type: 'number',
      format: 'double',
    },
    notes: {
      type: 'string',
    },
  },
} as const;
