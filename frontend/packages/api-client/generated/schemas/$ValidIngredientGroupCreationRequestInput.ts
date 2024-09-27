/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $ValidIngredientGroupCreationRequestInput = {
  properties: {
    description: {
      type: 'string',
    },
    members: {
      type: 'array',
      contains: {
        type: 'ValidIngredientGroupMemberCreationRequestInput',
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
