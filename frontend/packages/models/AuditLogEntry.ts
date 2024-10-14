// GENERATED CODE, DO NOT EDIT MANUALLY

import { ChangeLog } from './ChangeLog';

export interface IAuditLogEntry {
  belongsToHousehold: string;
  belongsToUser: string;
  changes: ChangeLog;
  createdAt: string;
  eventType: string;
  id: string;
  relevantID: string;
  resourceType: string;
}

export class AuditLogEntry implements IAuditLogEntry {
  belongsToHousehold: string;
  belongsToUser: string;
  changes: ChangeLog;
  createdAt: string;
  eventType: string;
  id: string;
  relevantID: string;
  resourceType: string;
  constructor(input: Partial<AuditLogEntry> = {}) {
    this.belongsToHousehold = input.belongsToHousehold || '';
    this.belongsToUser = input.belongsToUser || '';
    this.changes = input.changes || new ChangeLog();
    this.createdAt = input.createdAt || '';
    this.eventType = input.eventType || '';
    this.id = input.id || '';
    this.relevantID = input.relevantID || '';
    this.resourceType = input.resourceType || '';
  }
}
