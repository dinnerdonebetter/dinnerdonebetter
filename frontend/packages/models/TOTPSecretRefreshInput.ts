// GENERATED CODE, DO NOT EDIT MANUALLY

export interface ITOTPSecretRefreshInput {
  totpToken: string;
  currentPassword: string;
}

export class TOTPSecretRefreshInput implements ITOTPSecretRefreshInput {
  totpToken: string;
  currentPassword: string;
  constructor(input: Partial<TOTPSecretRefreshInput> = {}) {
    this.totpToken = input.totpToken || '';
    this.currentPassword = input.currentPassword || '';
  }
}
