// GENERATED CODE, DO NOT EDIT MANUALLY

import { Household } from './Household';
import { User } from './User';

export interface IHouseholdInvitation {
  archivedAt?: string;
  expiresAt: string;
  createdAt: string;
  id: string;
  toName: string;
  status: string;
  token: string;
  destinationHousehold: Household;
  fromUser: User;
  lastUpdatedAt?: string;
  note: string;
  statusNote: string;
  toEmail: string;
  toUser?: string;
}

export class HouseholdInvitation implements IHouseholdInvitation {
  archivedAt?: string;
  expiresAt: string;
  createdAt: string;
  id: string;
  toName: string;
  status: string;
  token: string;
  destinationHousehold: Household;
  fromUser: User;
  lastUpdatedAt?: string;
  note: string;
  statusNote: string;
  toEmail: string;
  toUser?: string;
  constructor(input: Partial<HouseholdInvitation> = {}) {
    this.archivedAt = input.archivedAt;
    this.expiresAt = input.expiresAt = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.toName = input.toName = '';
    this.status = input.status = '';
    this.token = input.token = '';
    this.destinationHousehold = input.destinationHousehold = new Household();
    this.fromUser = input.fromUser = new User();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.note = input.note = '';
    this.statusNote = input.statusNote = '';
    this.toEmail = input.toEmail = '';
    this.toUser = input.toUser;
  }
}
