// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredient } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateValidIngredientResponseConfig extends ResponseConfig<ValidIngredient> {
  constructor(status: number = 201, body?: ValidIngredient) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateValidIngredient = (resCfg: MockCreateValidIngredientResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredients`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
