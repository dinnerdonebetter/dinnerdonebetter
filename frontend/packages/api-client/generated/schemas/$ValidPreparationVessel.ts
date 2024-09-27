/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $ValidPreparationVessel = {
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
    instrument: {
      type: 'ValidVessel',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    notes: {
      type: 'string',
    },
    preparation: {
      type: 'ValidPreparation',
    },
  },
} as const;
