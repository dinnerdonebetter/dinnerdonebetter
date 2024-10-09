// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IPasswordUpdateInput {
  newPassword: string;
  totpToken: string;
  currentPassword: string;
}

export class PasswordUpdateInput implements IPasswordUpdateInput {
  newPassword: string;
  totpToken: string;
  currentPassword: string;
  constructor(input: Partial<PasswordUpdateInput> = {}) {
    this.newPassword = input.newPassword = '';
    this.totpToken = input.totpToken = '';
    this.currentPassword = input.currentPassword = '';
  }
}
