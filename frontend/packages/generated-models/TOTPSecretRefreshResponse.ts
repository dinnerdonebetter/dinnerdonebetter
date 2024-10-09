// GENERATED CODE, DO NOT EDIT MANUALLY



export interface ITOTPSecretRefreshResponse {
   qrCode: string;
 twoFactorSecret: string;

}

export class TOTPSecretRefreshResponse implements ITOTPSecretRefreshResponse {
   qrCode: string;
 twoFactorSecret: string;
constructor(input: Partial<TOTPSecretRefreshResponse> = {}) {
	 this.qrCode = input.qrCode = '';
 this.twoFactorSecret = input.twoFactorSecret = '';
}
}