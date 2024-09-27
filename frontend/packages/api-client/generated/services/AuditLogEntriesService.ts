/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { AuditLogEntry } from '../models/AuditLogEntry';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class AuditLogEntriesService {
  /**
   * Operation for fetching AuditLogEntry
   * @param auditLogEntryId
   * @returns any
   * @throws ApiError
   */
  public static getAuditLogEntryById(auditLogEntryId: string): CancelablePromise<
    APIResponse & {
      data?: AuditLogEntry;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/audit_log_entries/{auditLogEntryID}',
      path: {
        auditLogEntryID: auditLogEntryId,
      },
    });
  }
  /**
   * Operation for fetching AuditLogEntry
   * @returns any
   * @throws ApiError
   */
  public static getAuditLogEntriesForHousehold(): CancelablePromise<
    APIResponse & {
      data?: AuditLogEntry;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/audit_log_entries/for_household',
    });
  }
  /**
   * Operation for fetching AuditLogEntry
   * @returns any
   * @throws ApiError
   */
  public static getAuditLogEntriesForUser(): CancelablePromise<
    APIResponse & {
      data?: AuditLogEntry;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/audit_log_entries/for_user',
    });
  }
}
