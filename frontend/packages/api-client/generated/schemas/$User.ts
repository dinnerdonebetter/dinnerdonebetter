/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $User = {
  properties: {
    accountStatus: {
      type: 'string',
    },
    accountStatusExplanation: {
      type: 'string',
    },
    archivedAt: {
      type: 'string',
      format: 'date-time',
    },
    avatar: {
      type: 'string',
    },
    birthday: {
      type: 'string',
      format: 'date-time',
    },
    createdAt: {
      type: 'string',
      format: 'date-time',
    },
    emailAddress: {
      type: 'string',
      format: 'email',
    },
    emailAddressVerifiedAt: {
      type: 'string',
      format: 'date-time',
    },
    firstName: {
      type: 'string',
    },
    id: {
      type: 'string',
    },
    lastAcceptedPrivacyPolicy: {
      type: 'string',
      format: 'date-time',
    },
    lastAcceptedTOS: {
      type: 'string',
      format: 'date-time',
    },
    lastName: {
      type: 'string',
    },
    lastUpdatedAt: {
      type: 'string',
      format: 'date-time',
    },
    passwordLastChangedAt: {
      type: 'string',
      format: 'date-time',
    },
    requiresPasswordChange: {
      type: 'boolean',
    },
    serviceRoles: {
      type: 'string',
    },
    twoFactorSecretVerifiedAt: {
      type: 'string',
      format: 'date-time',
    },
    username: {
      type: 'string',
    },
  },
} as const;
