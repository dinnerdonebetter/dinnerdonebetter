// GENERATED CODE, DO NOT EDIT MANUALLY

import { Household } from './Household';
import { User } from './User';

export interface IHouseholdInvitation {
  lastUpdatedAt: string;
  status: string;
  expiresAt: string;
  statusNote: string;
  archivedAt: string;
  createdAt: string;
  note: string;
  toName: string;
  token: string;
  destinationHousehold: Household;
  id: string;
  toUser: string;
  fromUser: User;
  toEmail: string;
}

export class HouseholdInvitation implements IHouseholdInvitation {
  lastUpdatedAt: string;
  status: string;
  expiresAt: string;
  statusNote: string;
  archivedAt: string;
  createdAt: string;
  note: string;
  toName: string;
  token: string;
  destinationHousehold: Household;
  id: string;
  toUser: string;
  fromUser: User;
  toEmail: string;
  constructor(input: Partial<HouseholdInvitation> = {}) {
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.status = input.status || '';
    this.expiresAt = input.expiresAt || '';
    this.statusNote = input.statusNote || '';
    this.archivedAt = input.archivedAt || '';
    this.createdAt = input.createdAt || '';
    this.note = input.note || '';
    this.toName = input.toName || '';
    this.token = input.token || '';
    this.destinationHousehold = input.destinationHousehold || new Household();
    this.id = input.id || '';
    this.toUser = input.toUser || '';
    this.fromUser = input.fromUser || new User();
    this.toEmail = input.toEmail || '';
  }
}
