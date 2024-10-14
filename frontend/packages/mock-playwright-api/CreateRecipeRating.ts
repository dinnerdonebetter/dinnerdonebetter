// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeRating } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateRecipeRatingResponseConfig extends ResponseConfig<RecipeRating> {
  recipeID: string;

  constructor(recipeID: string, status: number = 201, body?: RecipeRating) {
    super();

    this.recipeID = recipeID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateRecipeRating = (resCfg: MockCreateRecipeRatingResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/ratings`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
