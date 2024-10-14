// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStep, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetRecipeStepsResponseConfig extends ResponseConfig<QueryFilteredResult<RecipeStep>> {
  recipeID: string;

  constructor(recipeID: string, status: number = 200, body: RecipeStep[] = []) {
    super();

    this.recipeID = recipeID;

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetRecipeStepss = (resCfg: MockGetRecipeStepsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps`,
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
