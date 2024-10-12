// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserLoginInput {
   password: string;
 totpToken: string;
 username: string;

}

export class UserLoginInput implements IUserLoginInput {
   password: string;
 totpToken: string;
 username: string;
constructor(input: Partial<UserLoginInput> = {}) {
	 this.password = input.password || '';
 this.totpToken = input.totpToken || '';
 this.username = input.username || '';
}
}