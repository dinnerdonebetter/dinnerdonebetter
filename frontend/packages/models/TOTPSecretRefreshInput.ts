// GENERATED CODE, DO NOT EDIT MANUALLY



export interface ITOTPSecretRefreshInput {
   currentPassword: string;
 totpToken: string;

}

export class TOTPSecretRefreshInput implements ITOTPSecretRefreshInput {
   currentPassword: string;
 totpToken: string;
constructor(input: Partial<TOTPSecretRefreshInput> = {}) {
	 this.currentPassword = input.currentPassword || '';
 this.totpToken = input.totpToken || '';
}
}