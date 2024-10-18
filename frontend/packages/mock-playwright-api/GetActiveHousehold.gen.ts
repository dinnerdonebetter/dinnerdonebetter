// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Household } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetActiveHouseholdResponseConfig extends ResponseConfig<Household> {
  constructor(status: number = 200, body?: Household) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetActiveHousehold = (resCfg: MockGetActiveHouseholdResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/households/current`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
