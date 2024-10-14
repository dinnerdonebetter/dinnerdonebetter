// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredient, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetValidIngredientsByPreparationResponseConfig extends ResponseConfig<
  QueryFilteredResult<ValidIngredient>
> {
  q: string;
  validPreparationID: string;

  constructor(q: string, validPreparationID: string, status: number = 200, body: ValidIngredient[] = []) {
    super();

    this.q = q;
    this.validPreparationID = validPreparationID;

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetValidIngredientsByPreparations = (resCfg: MockGetValidIngredientsByPreparationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredients/by_preparation/${resCfg.validPreparationID}`,
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
