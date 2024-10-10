// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidVessel, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetValidVesselsResponseConfig extends ResponseConfig<QueryFilteredResult<ValidVessel>> {
  constructor(status: number = 200, body: ValidVessel[] = []) {
    super();

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetValidVesselss = (resCfg: MockGetValidVesselsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_vessels`,
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
