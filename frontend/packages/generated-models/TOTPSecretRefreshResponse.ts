// GENERATED CODE, DO NOT EDIT MANUALLY

export interface ITOTPSecretRefreshResponse {
  twoFactorSecret: string;
  qrCode: string;
}

export class TOTPSecretRefreshResponse implements ITOTPSecretRefreshResponse {
  twoFactorSecret: string;
  qrCode: string;
  constructor(input: Partial<TOTPSecretRefreshResponse> = {}) {
    this.twoFactorSecret = input.twoFactorSecret = '';
    this.qrCode = input.qrCode = '';
  }
}
