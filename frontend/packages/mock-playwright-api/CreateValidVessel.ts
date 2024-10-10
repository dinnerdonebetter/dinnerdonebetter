// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidVessel } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateValidVesselResponseConfig extends ResponseConfig<ValidVessel> {
  constructor(status: number = 201, body?: ValidVessel) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateValidVessel = (resCfg: MockCreateValidVesselResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_vessels`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
