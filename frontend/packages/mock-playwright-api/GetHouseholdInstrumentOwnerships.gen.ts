// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { HouseholdInstrumentOwnership, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetHouseholdInstrumentOwnershipsResponseConfig extends ResponseConfig<
  QueryFilteredResult<HouseholdInstrumentOwnership>
> {
  constructor(status: number = 200, body: HouseholdInstrumentOwnership[] = []) {
    super();

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetHouseholdInstrumentOwnershipss = (resCfg: MockGetHouseholdInstrumentOwnershipsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/households/instruments`,
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
