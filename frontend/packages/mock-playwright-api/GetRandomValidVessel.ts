// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidVessel } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetRandomValidVesselResponseConfig extends ResponseConfig<ValidVessel> {
  constructor(status: number = 200, body?: ValidVessel) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetRandomValidVessel = (resCfg: MockGetRandomValidVesselResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_vessels/random`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
