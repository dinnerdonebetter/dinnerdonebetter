// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientStateIngredient, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetValidIngredientStateIngredientsByIngredientResponseConfig extends ResponseConfig<
  QueryFilteredResult<ValidIngredientStateIngredient>
> {
  validIngredientID: string;

  constructor(validIngredientID: string, status: number = 200, body: ValidIngredientStateIngredient[] = []) {
    super();

    this.validIngredientID = validIngredientID;

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetValidIngredientStateIngredientsByIngredients = (
  resCfg: MockGetValidIngredientStateIngredientsByIngredientResponseConfig,
) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_state_ingredients/by_ingredient/${resCfg.validIngredientID}`,
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
