// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidMeasurementUnitConversion } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockValidMeasurementUnitConversionsToUnitResponseConfig extends ResponseConfig<ValidMeasurementUnitConversion> {
  validMeasurementUnitID: string;

  constructor(validMeasurementUnitID: string, status: number = 200, body?: ValidMeasurementUnitConversion) {
    super();

    this.validMeasurementUnitID = validMeasurementUnitID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockValidMeasurementUnitConversionsToUnit = (
  resCfg: MockValidMeasurementUnitConversionsToUnitResponseConfig,
) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_measurement_conversions/to_unit/${resCfg.validMeasurementUnitID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
