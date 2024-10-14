// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientGroup } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateValidIngredientGroupResponseConfig extends ResponseConfig<ValidIngredientGroup> {
  constructor(status: number = 201, body?: ValidIngredientGroup) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateValidIngredientGroup = (resCfg: MockCreateValidIngredientGroupResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_groups`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
