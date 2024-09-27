/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ChangeLog } from './ChangeLog';
export type AuditLogEntry = {
  belongsToHousehold?: string;
  belongsToUser?: string;
  changes?: ChangeLog;
  createdAt?: string;
  eventType?: string;
  id?: string;
  relevantID?: string;
  resourceType?: string;
};
