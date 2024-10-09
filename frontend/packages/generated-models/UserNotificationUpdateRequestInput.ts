// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserNotificationUpdateRequestInput {
   status?: string;

}

export class UserNotificationUpdateRequestInput implements IUserNotificationUpdateRequestInput {
   status?: string;
constructor(input: Partial<UserNotificationUpdateRequestInput> = {}) {
	 this.status = input.status;
}
}