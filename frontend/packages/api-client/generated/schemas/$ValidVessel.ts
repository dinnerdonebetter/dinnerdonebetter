/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $ValidVessel = {
  properties: {
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    capacity: {
      type: 'number',
      format: 'double',
    },
    capacityUnit: {
      type: 'ValidMeasurementUnit',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    description: {
      type: 'string',
    },
    displayInSummaryLists: {
      type: 'boolean',
    },
    heightInMillimeters: {
      type: 'number',
      format: 'double',
    },
    iconPath: {
      type: 'string',
    },
    id: {
      type: 'string',
    },
    includeInGeneratedInstructions: {
      type: 'boolean',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    lengthInMillimeters: {
      type: 'number',
      format: 'double',
    },
    name: {
      type: 'string',
    },
    pluralName: {
      type: 'string',
    },
    shape: {
      type: 'string',
    },
    slug: {
      type: 'string',
    },
    usableForStorage: {
      type: 'boolean',
    },
    widthInMillimeters: {
      type: 'number',
      format: 'double',
    },
  },
} as const;
