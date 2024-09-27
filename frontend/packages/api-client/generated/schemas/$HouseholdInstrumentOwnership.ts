/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $HouseholdInstrumentOwnership = {
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
    id: {
      type: 'string',
    },
    instrument: {
      type: 'ValidInstrument',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    notes: {
      type: 'string',
    },
    quantity: {
      type: 'number',
      format: 'int64',
    },
  },
} as const;
