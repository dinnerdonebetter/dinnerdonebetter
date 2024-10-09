// GENERATED CODE, DO NOT EDIT MANUALLY

import { Household } from './Household';
import { User } from './User';

export interface IHouseholdInvitation {
  destinationHousehold: Household;
  archivedAt?: string;
  note: string;
  toEmail: string;
  toName: string;
  createdAt: string;
  id: string;
  statusNote: string;
  toUser?: string;
  token: string;
  fromUser: User;
  lastUpdatedAt?: string;
  status: string;
  expiresAt: string;
}

export class HouseholdInvitation implements IHouseholdInvitation {
  destinationHousehold: Household;
  archivedAt?: string;
  note: string;
  toEmail: string;
  toName: string;
  createdAt: string;
  id: string;
  statusNote: string;
  toUser?: string;
  token: string;
  fromUser: User;
  lastUpdatedAt?: string;
  status: string;
  expiresAt: string;
  constructor(input: Partial<HouseholdInvitation> = {}) {
    this.destinationHousehold = input.destinationHousehold = new Household();
    this.archivedAt = input.archivedAt;
    this.note = input.note = '';
    this.toEmail = input.toEmail = '';
    this.toName = input.toName = '';
    this.createdAt = input.createdAt = '';
    this.id = input.id = '';
    this.statusNote = input.statusNote = '';
    this.toUser = input.toUser;
    this.token = input.token = '';
    this.fromUser = input.fromUser = new User();
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.status = input.status = '';
    this.expiresAt = input.expiresAt = '';
  }
}
