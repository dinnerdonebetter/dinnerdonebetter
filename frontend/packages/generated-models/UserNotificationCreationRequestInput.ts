// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserNotificationCreationRequestInput {
   status: string;
 belongsToUser: string;
 content: string;

}

export class UserNotificationCreationRequestInput implements IUserNotificationCreationRequestInput {
   status: string;
 belongsToUser: string;
 content: string;
constructor(input: Partial<UserNotificationCreationRequestInput> = {}) {
	 this.status = input.status = '';
 this.belongsToUser = input.belongsToUser = '';
 this.content = input.content = '';
}
}