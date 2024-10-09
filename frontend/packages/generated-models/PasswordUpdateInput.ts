// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IPasswordUpdateInput {
  currentPassword: string;
  newPassword: string;
  totpToken: string;
}

export class PasswordUpdateInput implements IPasswordUpdateInput {
  currentPassword: string;
  newPassword: string;
  totpToken: string;
  constructor(input: Partial<PasswordUpdateInput> = {}) {
    this.currentPassword = input.currentPassword = '';
    this.newPassword = input.newPassword = '';
    this.totpToken = input.totpToken = '';
  }
}
