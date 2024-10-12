// GENERATED CODE, DO NOT EDIT MANUALLY



export interface ITOTPSecretVerificationInput {
   totpToken: string;
 userID: string;

}

export class TOTPSecretVerificationInput implements ITOTPSecretVerificationInput {
   totpToken: string;
 userID: string;
constructor(input: Partial<TOTPSecretVerificationInput> = {}) {
	 this.totpToken = input.totpToken || '';
 this.userID = input.userID || '';
}
}