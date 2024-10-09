// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserNotification {
  createdAt: string;
  id: string;
  lastUpdatedAt?: string;
  status: string;
  belongsToUser: string;
  content: string;
}

export class UserNotification implements IUserNotification {
  createdAt: string;
  id: string;
  lastUpdatedAt?: string;
  status: string;
  belongsToUser: string;
  content: string;
  constructor(input: Partial<UserNotification> = {}) {
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.status = input.status = '';
    this.belongsToUser = input.belongsToUser = '';
    this.content = input.content = '';
  }
}
