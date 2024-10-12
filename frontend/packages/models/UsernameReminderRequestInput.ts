// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUsernameReminderRequestInput {
   emailAddress: string;

}

export class UsernameReminderRequestInput implements IUsernameReminderRequestInput {
   emailAddress: string;
constructor(input: Partial<UsernameReminderRequestInput> = {}) {
	 this.emailAddress = input.emailAddress || '';
}
}