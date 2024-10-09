// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserDetailsUpdateRequestInput {
   currentPassword: string;
 firstName: string;
 lastName: string;
 totpToken: string;
 birthday: string;

}

export class UserDetailsUpdateRequestInput implements IUserDetailsUpdateRequestInput {
   currentPassword: string;
 firstName: string;
 lastName: string;
 totpToken: string;
 birthday: string;
constructor(input: Partial<UserDetailsUpdateRequestInput> = {}) {
	 this.currentPassword = input.currentPassword = '';
 this.firstName = input.firstName = '';
 this.lastName = input.lastName = '';
 this.totpToken = input.totpToken = '';
 this.birthday = input.birthday = '';
}
}