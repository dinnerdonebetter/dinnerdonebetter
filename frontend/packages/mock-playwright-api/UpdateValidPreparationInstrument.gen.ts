// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidPreparationInstrument } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdateValidPreparationInstrumentResponseConfig extends ResponseConfig<ValidPreparationInstrument> {
  validPreparationVesselID: string;

  constructor(validPreparationVesselID: string, status: number = 200, body?: ValidPreparationInstrument) {
    super();

    this.validPreparationVesselID = validPreparationVesselID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdateValidPreparationInstrument = (resCfg: MockUpdateValidPreparationInstrumentResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparation_instruments/${resCfg.validPreparationVesselID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
