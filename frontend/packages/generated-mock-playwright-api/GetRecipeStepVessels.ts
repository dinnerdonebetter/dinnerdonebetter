// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStepVessel, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetRecipeStepVesselsResponseConfig extends ResponseConfig<QueryFilteredResult<RecipeStepVessel>> {
  recipeID: string;
  recipeStepID: string;

  constructor(recipeID: string, recipeStepID: string, status: number = 200, body: RecipeStepVessel[] = []) {
    super();

    this.recipeID = recipeID;
    this.recipeStepID = recipeStepID;

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetRecipeStepVesselss = (resCfg: MockGetRecipeStepVesselsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps/${resCfg.recipeStepID}/vessels`,
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
