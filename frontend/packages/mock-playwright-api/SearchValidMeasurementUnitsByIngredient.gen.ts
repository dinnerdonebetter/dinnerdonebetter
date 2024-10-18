// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidMeasurementUnit, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockSearchValidMeasurementUnitsByIngredientResponseConfig extends ResponseConfig<
  QueryFilteredResult<ValidMeasurementUnit>
> {
  q: string;
  validIngredientID: string;

  constructor(q: string, validIngredientID: string, status: number = 200, body: ValidMeasurementUnit[] = []) {
    super();

    this.q = q;
    this.validIngredientID = validIngredientID;

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockSearchValidMeasurementUnitsByIngredients = (
  resCfg: MockSearchValidMeasurementUnitsByIngredientResponseConfig,
) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_measurement_units/by_ingredient/${resCfg.validIngredientID}`,
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
