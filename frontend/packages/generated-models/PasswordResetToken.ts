// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IPasswordResetToken {
  token: string;
  archivedAt?: string;
  belongsToUser: string;
  createdAt: string;
  expiresAt: string;
  id: string;
  lastUpdatedAt?: string;
}

export class PasswordResetToken implements IPasswordResetToken {
  token: string;
  archivedAt?: string;
  belongsToUser: string;
  createdAt: string;
  expiresAt: string;
  id: string;
  lastUpdatedAt?: string;
  constructor(input: Partial<PasswordResetToken> = {}) {
    this.token = input.token = '';
    this.archivedAt = input.archivedAt;
    this.belongsToUser = input.belongsToUser = '';
    this.createdAt = input.createdAt = '';
    this.expiresAt = input.expiresAt = '';
    this.id = input.id = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
  }
}
