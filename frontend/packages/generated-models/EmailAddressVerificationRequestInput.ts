// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IEmailAddressVerificationRequestInput {
   emailVerificationToken: string;

}

export class EmailAddressVerificationRequestInput implements IEmailAddressVerificationRequestInput {
   emailVerificationToken: string;
constructor(input: Partial<EmailAddressVerificationRequestInput> = {}) {
	 this.emailVerificationToken = input.emailVerificationToken = '';
}
}