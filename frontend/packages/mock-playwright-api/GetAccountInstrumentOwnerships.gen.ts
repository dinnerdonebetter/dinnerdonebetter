// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { AccountInstrumentOwnership, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetAccountInstrumentOwnershipsResponseConfig extends ResponseConfig<
  QueryFilteredResult<AccountInstrumentOwnership>
> {
  constructor(status: number = 200, body: AccountInstrumentOwnership[] = []) {
    super();

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetAccountInstrumentOwnershipss = (resCfg: MockGetAccountInstrumentOwnershipsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/accounts/instruments`,
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
