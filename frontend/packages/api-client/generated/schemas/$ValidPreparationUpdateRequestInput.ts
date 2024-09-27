/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $ValidPreparationUpdateRequestInput = {
  properties: {
    conditionExpressionRequired: {
      type: 'boolean',
    },
    consumesVessel: {
      type: 'boolean',
    },
    description: {
      type: 'string',
    },
    iconPath: {
      type: 'string',
    },
    maximumIngredientCount: {
      type: 'number',
      format: 'int32',
    },
    maximumInstrumentCount: {
      type: 'number',
      format: 'int32',
    },
    maximumVesselCount: {
      type: 'number',
      format: 'int32',
    },
    minimumIngredientCount: {
      type: 'number',
      format: 'int32',
    },
    minimumInstrumentCount: {
      type: 'number',
      format: 'int32',
    },
    minimumVesselCount: {
      type: 'number',
      format: 'int32',
    },
    name: {
      type: 'string',
    },
    onlyForVessels: {
      type: 'boolean',
    },
    pastTense: {
      type: 'string',
    },
    restrictToIngredients: {
      type: 'boolean',
    },
    slug: {
      type: 'string',
    },
    temperatureRequired: {
      type: 'boolean',
    },
    timeEstimateRequired: {
      type: 'boolean',
    },
    yieldsNothing: {
      type: 'boolean',
    },
  },
} as const;
