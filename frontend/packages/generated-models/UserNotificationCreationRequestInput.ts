// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserNotificationCreationRequestInput {
  content: string;
  status: string;
  belongsToUser: string;
}

export class UserNotificationCreationRequestInput implements IUserNotificationCreationRequestInput {
  content: string;
  status: string;
  belongsToUser: string;
  constructor(input: Partial<UserNotificationCreationRequestInput> = {}) {
    this.content = input.content = '';
    this.status = input.status = '';
    this.belongsToUser = input.belongsToUser = '';
  }
}
