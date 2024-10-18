// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidPreparationVessel } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetValidPreparationVesselResponseConfig extends ResponseConfig<ValidPreparationVessel> {
  validPreparationVesselID: string;

  constructor(validPreparationVesselID: string, status: number = 200, body?: ValidPreparationVessel) {
    super();

    this.validPreparationVesselID = validPreparationVesselID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetValidPreparationVessel = (resCfg: MockGetValidPreparationVesselResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparation_vessels/${resCfg.validPreparationVesselID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
