// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Recipe } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateRecipeResponseConfig extends ResponseConfig<Recipe> {
  constructor(status: number = 201, body?: Recipe) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateRecipe = (resCfg: MockCreateRecipeResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
