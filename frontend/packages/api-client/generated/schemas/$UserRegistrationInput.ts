/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export const $UserRegistrationInput = {
  properties: {
    acceptedPrivacyPolicy: {
      type: 'boolean',
    },
    acceptedTOS: {
      type: 'boolean',
    },
    birthday: {
      type: 'string',
      format: 'date-time',
    },
    emailAddress: {
      type: 'string',
      format: 'email',
    },
    firstName: {
      type: 'string',
    },
    householdName: {
      type: 'string',
    },
    invitationID: {
      type: 'string',
    },
    invitationToken: {
      type: 'string',
    },
    lastName: {
      type: 'string',
    },
    password: {
      type: 'string',
      format: 'password',
    },
    username: {
      type: 'string',
    },
  },
} as const;
