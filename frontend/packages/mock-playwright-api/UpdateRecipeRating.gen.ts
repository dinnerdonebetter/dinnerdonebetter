// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeRating } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdateRecipeRatingResponseConfig extends ResponseConfig<RecipeRating> {
  recipeID: string;
  recipeRatingID: string;

  constructor(recipeID: string, recipeRatingID: string, status: number = 200, body?: RecipeRating) {
    super();

    this.recipeID = recipeID;
    this.recipeRatingID = recipeRatingID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdateRecipeRating = (resCfg: MockUpdateRecipeRatingResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/ratings/${resCfg.recipeRatingID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
