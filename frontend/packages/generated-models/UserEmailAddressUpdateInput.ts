// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserEmailAddressUpdateInput {
  currentPassword: string;
  newEmailAddress: string;
  totpToken: string;
}

export class UserEmailAddressUpdateInput implements IUserEmailAddressUpdateInput {
  currentPassword: string;
  newEmailAddress: string;
  totpToken: string;
  constructor(input: Partial<UserEmailAddressUpdateInput> = {}) {
    this.currentPassword = input.currentPassword = '';
    this.newEmailAddress = input.newEmailAddress = '';
    this.totpToken = input.totpToken = '';
  }
}
