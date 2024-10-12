// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { AuditLogEntry } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetAuditLogEntryByIDResponseConfig extends ResponseConfig<AuditLogEntry> {
		   auditLogEntryID: string;
		

		  constructor( auditLogEntryID: string, status: number = 200, body?: AuditLogEntry) {
		    super();

		 this.auditLogEntryID = auditLogEntryID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetAuditLogEntryByID = (resCfg: MockGetAuditLogEntryByIDResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/audit_log_entries/${resCfg.auditLogEntryID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};