// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientState, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetValidIngredientStatesResponseConfig extends ResponseConfig<
  QueryFilteredResult<ValidIngredientState>
> {
  constructor(status: number = 200, body: ValidIngredientState[] = []) {
    super();

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetValidIngredientStatess = (resCfg: MockGetValidIngredientStatesResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_states`,
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
