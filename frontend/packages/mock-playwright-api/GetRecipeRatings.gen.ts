// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeRating, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetRecipeRatingsResponseConfig extends ResponseConfig<QueryFilteredResult<RecipeRating>> {
  recipeID: string;

  constructor(recipeID: string, status: number = 200, body: RecipeRating[] = []) {
    super();

    this.recipeID = recipeID;

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetRecipeRatingss = (resCfg: MockGetRecipeRatingsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/ratings`,
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
