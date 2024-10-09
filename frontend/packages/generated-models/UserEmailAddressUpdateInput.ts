// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserEmailAddressUpdateInput {
  totpToken: string;
  currentPassword: string;
  newEmailAddress: string;
}

export class UserEmailAddressUpdateInput implements IUserEmailAddressUpdateInput {
  totpToken: string;
  currentPassword: string;
  newEmailAddress: string;
  constructor(input: Partial<UserEmailAddressUpdateInput> = {}) {
    this.totpToken = input.totpToken = '';
    this.currentPassword = input.currentPassword = '';
    this.newEmailAddress = input.newEmailAddress = '';
  }
}
