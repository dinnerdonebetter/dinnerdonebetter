// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserEmailAddressUpdateInput {
   newEmailAddress: string;
 totpToken: string;
 currentPassword: string;

}

export class UserEmailAddressUpdateInput implements IUserEmailAddressUpdateInput {
   newEmailAddress: string;
 totpToken: string;
 currentPassword: string;
constructor(input: Partial<UserEmailAddressUpdateInput> = {}) {
	 this.newEmailAddress = input.newEmailAddress = '';
 this.totpToken = input.totpToken = '';
 this.currentPassword = input.currentPassword = '';
}
}