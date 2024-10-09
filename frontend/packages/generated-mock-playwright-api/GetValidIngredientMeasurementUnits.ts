// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientMeasurementUnit, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetValidIngredientMeasurementUnitsResponseConfig extends ResponseConfig<
  QueryFilteredResult<ValidIngredientMeasurementUnit>
> {
  constructor(status: number = 200, body: ValidIngredientMeasurementUnit[] = []) {
    super();

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetValidIngredientMeasurementUnitss = (
  resCfg: MockGetValidIngredientMeasurementUnitsResponseConfig,
) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_measurement_units`,
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
