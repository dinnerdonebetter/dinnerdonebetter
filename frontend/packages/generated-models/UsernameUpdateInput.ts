// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUsernameUpdateInput {
  currentPassword: string;
  newUsername: string;
  totpToken: string;
}

export class UsernameUpdateInput implements IUsernameUpdateInput {
  currentPassword: string;
  newUsername: string;
  totpToken: string;
  constructor(input: Partial<UsernameUpdateInput> = {}) {
    this.currentPassword = input.currentPassword = '';
    this.newUsername = input.newUsername = '';
    this.totpToken = input.totpToken = '';
  }
}
