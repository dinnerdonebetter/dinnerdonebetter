// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStep } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateRecipeStepResponseConfig extends ResponseConfig<RecipeStep> {
  recipeID: string;

  constructor(recipeID: string, status: number = 201, body?: RecipeStep) {
    super();

    this.recipeID = recipeID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateRecipeStep = (resCfg: MockCreateRecipeStepResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
