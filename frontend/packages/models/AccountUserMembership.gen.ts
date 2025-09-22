// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IAccountUserMembership {
  archivedAt: string;
  belongsToAccount: string;
  belongsToUser: string;
  createdAt: string;
  defaultAccount: boolean;
  accountRole: string;
  id: string;
  lastUpdatedAt: string;
}

export class AccountUserMembership implements IAccountUserMembership {
  archivedAt: string;
  belongsToAccount: string;
  belongsToUser: string;
  createdAt: string;
  defaultAccount: boolean;
  accountRole: string;
  id: string;
  lastUpdatedAt: string;
  constructor(input: Partial<AccountUserMembership> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToAccount = input.belongsToAccount || '';
    this.belongsToUser = input.belongsToUser || '';
    this.createdAt = input.createdAt || '';
    this.defaultAccount = input.defaultAccount || false;
    this.accountRole = input.accountRole || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
  }
}
