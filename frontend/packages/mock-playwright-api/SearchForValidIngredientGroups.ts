// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientGroup, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockSearchForValidIngredientGroupsResponseConfig extends ResponseConfig<
  QueryFilteredResult<ValidIngredientGroup>
> {
  q: string;

  constructor(q: string, status: number = 200, body: ValidIngredientGroup[] = []) {
    super();

    this.q = q;

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockSearchForValidIngredientGroupss = (resCfg: MockSearchForValidIngredientGroupsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_groups/search`,
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
