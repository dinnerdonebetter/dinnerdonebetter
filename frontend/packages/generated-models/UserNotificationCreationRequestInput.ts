// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserNotificationCreationRequestInput {
  belongsToUser: string;
  content: string;
  status: string;
}

export class UserNotificationCreationRequestInput implements IUserNotificationCreationRequestInput {
  belongsToUser: string;
  content: string;
  status: string;
  constructor(input: Partial<UserNotificationCreationRequestInput> = {}) {
    this.belongsToUser = input.belongsToUser = '';
    this.content = input.content = '';
    this.status = input.status = '';
  }
}
