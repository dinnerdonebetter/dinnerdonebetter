// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidVessel } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveValidVesselResponseConfig extends ResponseConfig<ValidVessel> {
  validVesselID: string;

  constructor(validVesselID: string, status: number = 202, body?: ValidVessel) {
    super();

    this.validVesselID = validVesselID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveValidVessel = (resCfg: MockArchiveValidVesselResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_vessels/${resCfg.validVesselID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
