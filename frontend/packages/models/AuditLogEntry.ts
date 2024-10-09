// GENERATED CODE, DO NOT EDIT MANUALLY

import { ChangeLog } from './ChangeLog';

export interface IAuditLogEntry {
  id: string;
  relevantID: string;
  resourceType: string;
  belongsToHousehold: string;
  belongsToUser: string;
  changes: ChangeLog;
  createdAt: string;
  eventType: string;
}

export class AuditLogEntry implements IAuditLogEntry {
  id: string;
  relevantID: string;
  resourceType: string;
  belongsToHousehold: string;
  belongsToUser: string;
  changes: ChangeLog;
  createdAt: string;
  eventType: string;
  constructor(input: Partial<AuditLogEntry> = {}) {
    this.id = input.id || '';
    this.relevantID = input.relevantID || '';
    this.resourceType = input.resourceType || '';
    this.belongsToHousehold = input.belongsToHousehold || '';
    this.belongsToUser = input.belongsToUser || '';
    this.changes = input.changes || new ChangeLog();
    this.createdAt = input.createdAt || '';
    this.eventType = input.eventType || '';
  }
}
