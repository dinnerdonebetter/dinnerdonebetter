// GENERATED CODE, DO NOT EDIT MANUALLY

import { Household } from './Household';
import { User } from './User';

export interface IHouseholdInvitation {
  statusNote: string;
  toEmail: string;
  archivedAt?: string;
  id: string;
  toName: string;
  token: string;
  createdAt: string;
  destinationHousehold: Household;
  note: string;
  expiresAt: string;
  fromUser: User;
  toUser?: string;
  lastUpdatedAt?: string;
  status: string;
}

export class HouseholdInvitation implements IHouseholdInvitation {
  statusNote: string;
  toEmail: string;
  archivedAt?: string;
  id: string;
  toName: string;
  token: string;
  createdAt: string;
  destinationHousehold: Household;
  note: string;
  expiresAt: string;
  fromUser: User;
  toUser?: string;
  lastUpdatedAt?: string;
  status: string;
  constructor(input: Partial<HouseholdInvitation> = {}) {
    this.statusNote = input.statusNote = '';
    this.toEmail = input.toEmail = '';
    this.archivedAt = input.archivedAt;
    this.id = input.id = '';
    this.toName = input.toName = '';
    this.token = input.token = '';
    this.createdAt = input.createdAt = '';
    this.destinationHousehold = input.destinationHousehold = new Household();
    this.note = input.note = '';
    this.expiresAt = input.expiresAt = '';
    this.fromUser = input.fromUser = new User();
    this.toUser = input.toUser;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.status = input.status = '';
  }
}
