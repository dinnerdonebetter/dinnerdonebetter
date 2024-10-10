// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IPasswordResetToken {
  archivedAt: string;
  belongsToUser: string;
  createdAt: string;
  expiresAt: string;
  id: string;
  lastUpdatedAt: string;
  token: string;
}

export class PasswordResetToken implements IPasswordResetToken {
  archivedAt: string;
  belongsToUser: string;
  createdAt: string;
  expiresAt: string;
  id: string;
  lastUpdatedAt: string;
  token: string;
  constructor(input: Partial<PasswordResetToken> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.belongsToUser = input.belongsToUser || '';
    this.createdAt = input.createdAt || '';
    this.expiresAt = input.expiresAt || '';
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.token = input.token || '';
  }
}
