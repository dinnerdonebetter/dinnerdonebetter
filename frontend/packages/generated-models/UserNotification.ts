// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserNotification {
   belongsToUser: string;
 content: string;
 createdAt: string;
 id: string;
 lastUpdatedAt?: string;
 status: string;

}

export class UserNotification implements IUserNotification {
   belongsToUser: string;
 content: string;
 createdAt: string;
 id: string;
 lastUpdatedAt?: string;
 status: string;
constructor(input: Partial<UserNotification> = {}) {
	 this.belongsToUser = input.belongsToUser = '';
 this.content = input.content = '';
 this.createdAt = input.createdAt = '';
 this.id = input.id = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.status = input.status = '';
}
}