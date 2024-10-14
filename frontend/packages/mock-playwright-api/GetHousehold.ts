// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Household } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetHouseholdResponseConfig extends ResponseConfig<Household> {
  householdID: string;

  constructor(householdID: string, status: number = 200, body?: Household) {
    super();

    this.householdID = householdID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetHousehold = (resCfg: MockGetHouseholdResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/households/${resCfg.householdID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
