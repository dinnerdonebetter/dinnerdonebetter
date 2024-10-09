// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUsernameUpdateInput {
  totpToken: string;
  currentPassword: string;
  newUsername: string;
}

export class UsernameUpdateInput implements IUsernameUpdateInput {
  totpToken: string;
  currentPassword: string;
  newUsername: string;
  constructor(input: Partial<UsernameUpdateInput> = {}) {
    this.totpToken = input.totpToken = '';
    this.currentPassword = input.currentPassword = '';
    this.newUsername = input.newUsername = '';
  }
}
