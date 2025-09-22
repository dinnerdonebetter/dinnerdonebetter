// GENERATED CODE, DO NOT EDIT MANUALLY

import { User } from './User.gen';

export interface IAccountUserMembershipWithUser {
  archivedAt: string;
  belongsToAccount: string;
  belongsToUser: User;
  createdAt: string;
  defaultAccount: boolean;
  accountRole: string;
  id: string;
  lastUpdatedAt: string;
}

export class AccountUserMembershipWithUser implements IAccountUserMembershipWithUser {
  archivedAt: string;
  belongsToAccount: string;
  belongsToUser: User;
  createdAt: string;
  defaultAccount: boolean;
  accountRole: string;
  id: string;
  lastUpdatedAt: string;
  constructor(input: Partial<AccountUserMembershipWithUser> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToAccount = input.belongsToAccount || '';
    this.belongsToUser = input.belongsToUser || new User();
    this.createdAt = input.createdAt || '';
    this.defaultAccount = input.defaultAccount || false;
    this.accountRole = input.accountRole || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
  }
}
