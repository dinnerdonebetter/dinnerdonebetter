// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserRegistrationInput {
  birthday?: string;
  lastName: string;
  password: string;
  username: string;
  acceptedPrivacyPolicy: boolean;
  acceptedTOS: boolean;
  emailAddress: string;
  firstName: string;
  householdName: string;
  invitationID: string;
  invitationToken: string;
}

export class UserRegistrationInput implements IUserRegistrationInput {
  birthday?: string;
  lastName: string;
  password: string;
  username: string;
  acceptedPrivacyPolicy: boolean;
  acceptedTOS: boolean;
  emailAddress: string;
  firstName: string;
  householdName: string;
  invitationID: string;
  invitationToken: string;
  constructor(input: Partial<UserRegistrationInput> = {}) {
    this.birthday = input.birthday;
    this.lastName = input.lastName = '';
    this.password = input.password = '';
    this.username = input.username = '';
    this.acceptedPrivacyPolicy = input.acceptedPrivacyPolicy = false;
    this.acceptedTOS = input.acceptedTOS = false;
    this.emailAddress = input.emailAddress = '';
    this.firstName = input.firstName = '';
    this.householdName = input.householdName = '';
    this.invitationID = input.invitationID = '';
    this.invitationToken = input.invitationToken = '';
  }
}
