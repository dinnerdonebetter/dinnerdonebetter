// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientMeasurementUnit } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdateValidIngredientMeasurementUnitResponseConfig extends ResponseConfig<ValidIngredientMeasurementUnit> {
  validIngredientMeasurementUnitID: string;

  constructor(validIngredientMeasurementUnitID: string, status: number = 200, body?: ValidIngredientMeasurementUnit) {
    super();

    this.validIngredientMeasurementUnitID = validIngredientMeasurementUnitID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdateValidIngredientMeasurementUnit = (
  resCfg: MockUpdateValidIngredientMeasurementUnitResponseConfig,
) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_measurement_units/${resCfg.validIngredientMeasurementUnitID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
