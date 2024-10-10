// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { AuditLogEntry } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetAuditLogEntriesForUserResponseConfig extends ResponseConfig<AuditLogEntry> {
  constructor(status: number = 200, body?: AuditLogEntry) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetAuditLogEntriesForUser = (resCfg: MockGetAuditLogEntriesForUserResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/audit_log_entries/for_user`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
