// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IPasswordResetTokenCreationRequestInput {
   emailAddress: string;

}

export class PasswordResetTokenCreationRequestInput implements IPasswordResetTokenCreationRequestInput {
   emailAddress: string;
constructor(input: Partial<PasswordResetTokenCreationRequestInput> = {}) {
	 this.emailAddress = input.emailAddress || '';
}
}