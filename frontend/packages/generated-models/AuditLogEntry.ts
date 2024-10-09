// GENERATED CODE, DO NOT EDIT MANUALLY

import { ChangeLog } from './ChangeLog';

export interface IAuditLogEntry {
  createdAt: string;
  eventType: string;
  id: string;
  relevantID: string;
  resourceType: string;
  belongsToHousehold?: string;
  belongsToUser: string;
  changes: ChangeLog;
}

export class AuditLogEntry implements IAuditLogEntry {
  createdAt: string;
  eventType: string;
  id: string;
  relevantID: string;
  resourceType: string;
  belongsToHousehold?: string;
  belongsToUser: string;
  changes: ChangeLog;
  constructor(input: Partial<AuditLogEntry> = {}) {
    this.createdAt = input.createdAt = '';
    this.eventType = input.eventType = '';
    this.id = input.id = '';
    this.relevantID = input.relevantID = '';
    this.resourceType = input.resourceType = '';
    this.belongsToHousehold = input.belongsToHousehold;
    this.belongsToUser = input.belongsToUser = '';
    this.changes = input.changes = new ChangeLog();
  }
}
