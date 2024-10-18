// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { AuditLogEntry, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetAuditLogEntriesForUserResponseConfig extends ResponseConfig<QueryFilteredResult<AuditLogEntry>> {
  constructor(status: number = 200, body: AuditLogEntry[] = []) {
    super();

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetAuditLogEntriesForUsers = (resCfg: MockGetAuditLogEntriesForUserResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/audit_log_entries/for_user`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        if (resCfg.body && resCfg.filter) resCfg.body.limit = resCfg.filter.limit;
        if (resCfg.body && resCfg.filter) resCfg.body.page = resCfg.filter.page;

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
