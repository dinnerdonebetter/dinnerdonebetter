// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserRegistrationInput {
  firstName: string;
  lastName: string;
  password: string;
  username: string;
  emailAddress: string;
  acceptedTOS: boolean;
  birthday?: string;
  householdName: string;
  invitationID: string;
  invitationToken: string;
  acceptedPrivacyPolicy: boolean;
}

export class UserRegistrationInput implements IUserRegistrationInput {
  firstName: string;
  lastName: string;
  password: string;
  username: string;
  emailAddress: string;
  acceptedTOS: boolean;
  birthday?: string;
  householdName: string;
  invitationID: string;
  invitationToken: string;
  acceptedPrivacyPolicy: boolean;
  constructor(input: Partial<UserRegistrationInput> = {}) {
    this.firstName = input.firstName = '';
    this.lastName = input.lastName = '';
    this.password = input.password = '';
    this.username = input.username = '';
    this.emailAddress = input.emailAddress = '';
    this.acceptedTOS = input.acceptedTOS = false;
    this.birthday = input.birthday;
    this.householdName = input.householdName = '';
    this.invitationID = input.invitationID = '';
    this.invitationToken = input.invitationToken = '';
    this.acceptedPrivacyPolicy = input.acceptedPrivacyPolicy = false;
  }
}
