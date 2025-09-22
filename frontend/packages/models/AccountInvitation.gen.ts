// GENERATED CODE, DO NOT EDIT MANUALLY

import { Account } from './Account.gen';
import { User } from './User.gen';

export interface IAccountInvitation {
  archivedAt: string;
  createdAt: string;
  destinationAccount: Account;
  expiresAt: string;
  fromUser: User;
  id: string;
  lastUpdatedAt: string;
  note: string;
  status: string;
  statusNote: string;
  toEmail: string;
  toName: string;
  toUser: string;
  token: string;
}

export class AccountInvitation implements IAccountInvitation {
  archivedAt: string;
  createdAt: string;
  destinationAccount: Account;
  expiresAt: string;
  fromUser: User;
  id: string;
  lastUpdatedAt: string;
  note: string;
  status: string;
  statusNote: string;
  toEmail: string;
  toName: string;
  toUser: string;
  token: string;
  constructor(input: Partial<AccountInvitation> = {}) {
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.destinationAccount = input.destinationAccount || new Account();
    this.expiresAt = input.expiresAt || '';
    this.fromUser = input.fromUser || new User();
    this.id = input.id || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.note = input.note || '';
    this.status = input.status || '';
    this.statusNote = input.statusNote || '';
    this.toEmail = input.toEmail || '';
    this.toName = input.toName || '';
    this.toUser = input.toUser || '';
    this.token = input.token || '';
  }
}
