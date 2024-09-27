/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $ValidMeasurementUnitConversion = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    from: {
      type: 'ValidMeasurementUnit',
    },
    id: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    modifier: {
      type: 'number',
      format: 'double',
    },
    notes: {
      type: 'string',
    },
    onlyForIngredient: {
      type: 'ValidIngredient',
    },
    to: {
      type: 'ValidMeasurementUnit',
    },
  },
} as const;
