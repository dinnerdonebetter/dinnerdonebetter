// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserDetailsUpdateRequestInput {
  birthday: string;
  currentPassword: string;
  firstName: string;
  lastName: string;
  totpToken: string;
}

export class UserDetailsUpdateRequestInput implements IUserDetailsUpdateRequestInput {
  birthday: string;
  currentPassword: string;
  firstName: string;
  lastName: string;
  totpToken: string;
  constructor(input: Partial<UserDetailsUpdateRequestInput> = {}) {
    this.birthday = input.birthday || '';
    this.currentPassword = input.currentPassword || '';
    this.firstName = input.firstName || '';
    this.lastName = input.lastName || '';
    this.totpToken = input.totpToken || '';
  }
}
