// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidMeasurementUnit } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveValidMeasurementUnitResponseConfig extends ResponseConfig<ValidMeasurementUnit> {
  validMeasurementUnitID: string;

  constructor(validMeasurementUnitID: string, status: number = 202, body?: ValidMeasurementUnit) {
    super();

    this.validMeasurementUnitID = validMeasurementUnitID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveValidMeasurementUnit = (resCfg: MockArchiveValidMeasurementUnitResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_measurement_units/${resCfg.validMeasurementUnitID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
