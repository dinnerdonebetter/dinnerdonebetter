// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserRegistrationInput {
  householdName: string;
  lastName: string;
  emailAddress: string;
  firstName: string;
  invitationID: string;
  invitationToken: string;
  password: string;
  acceptedPrivacyPolicy: boolean;
  acceptedTOS: boolean;
  birthday?: string;
  username: string;
}

export class UserRegistrationInput implements IUserRegistrationInput {
  householdName: string;
  lastName: string;
  emailAddress: string;
  firstName: string;
  invitationID: string;
  invitationToken: string;
  password: string;
  acceptedPrivacyPolicy: boolean;
  acceptedTOS: boolean;
  birthday?: string;
  username: string;
  constructor(input: Partial<UserRegistrationInput> = {}) {
    this.householdName = input.householdName = '';
    this.lastName = input.lastName = '';
    this.emailAddress = input.emailAddress = '';
    this.firstName = input.firstName = '';
    this.invitationID = input.invitationID = '';
    this.invitationToken = input.invitationToken = '';
    this.password = input.password = '';
    this.acceptedPrivacyPolicy = input.acceptedPrivacyPolicy = false;
    this.acceptedTOS = input.acceptedTOS = false;
    this.birthday = input.birthday;
    this.username = input.username = '';
  }
}
