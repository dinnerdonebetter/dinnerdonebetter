// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IPasswordResetTokenRedemptionRequestInput {
  newPassword: string;
  token: string;
}

export class PasswordResetTokenRedemptionRequestInput implements IPasswordResetTokenRedemptionRequestInput {
  newPassword: string;
  token: string;
  constructor(input: Partial<PasswordResetTokenRedemptionRequestInput> = {}) {
    this.newPassword = input.newPassword || '';
    this.token = input.token || '';
  }
}
