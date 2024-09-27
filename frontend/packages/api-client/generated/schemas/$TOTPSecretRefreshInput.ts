/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $TOTPSecretRefreshInput = {
  properties: {
    currentPassword: {
      type: 'string',
      format: 'password',
    },
    totpToken: {
      type: 'string',
    },
  },
} as const;
