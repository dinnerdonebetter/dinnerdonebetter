// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserDetailsUpdateRequestInput {
  firstName: string;
  lastName: string;
  totpToken: string;
  birthday: string;
  currentPassword: string;
}

export class UserDetailsUpdateRequestInput implements IUserDetailsUpdateRequestInput {
  firstName: string;
  lastName: string;
  totpToken: string;
  birthday: string;
  currentPassword: string;
  constructor(input: Partial<UserDetailsUpdateRequestInput> = {}) {
    this.firstName = input.firstName = '';
    this.lastName = input.lastName = '';
    this.totpToken = input.totpToken = '';
    this.birthday = input.birthday = '';
    this.currentPassword = input.currentPassword = '';
  }
}
