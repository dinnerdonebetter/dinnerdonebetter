// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientMeasurementUnit } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateValidIngredientMeasurementUnitResponseConfig extends ResponseConfig<ValidIngredientMeasurementUnit> {
  constructor(status: number = 201, body?: ValidIngredientMeasurementUnit) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateValidIngredientMeasurementUnit = (
  resCfg: MockCreateValidIngredientMeasurementUnitResponseConfig,
) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_measurement_units`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
