// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserRegistrationInput {
  invitationToken: string;
  lastName: string;
  username: string;
  emailAddress: string;
  firstName: string;
  householdName: string;
  invitationID: string;
  password: string;
  acceptedPrivacyPolicy: boolean;
  acceptedTOS: boolean;
  birthday: string;
}

export class UserRegistrationInput implements IUserRegistrationInput {
  invitationToken: string;
  lastName: string;
  username: string;
  emailAddress: string;
  firstName: string;
  householdName: string;
  invitationID: string;
  password: string;
  acceptedPrivacyPolicy: boolean;
  acceptedTOS: boolean;
  birthday: string;
  constructor(input: Partial<UserRegistrationInput> = {}) {
    this.invitationToken = input.invitationToken || '';
    this.lastName = input.lastName || '';
    this.username = input.username || '';
    this.emailAddress = input.emailAddress || '';
    this.firstName = input.firstName || '';
    this.householdName = input.householdName || '';
    this.invitationID = input.invitationID || '';
    this.password = input.password || '';
    this.acceptedPrivacyPolicy = input.acceptedPrivacyPolicy || false;
    this.acceptedTOS = input.acceptedTOS || false;
    this.birthday = input.birthday || '';
  }
}
