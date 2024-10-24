// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStep } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdateRecipeStepResponseConfig extends ResponseConfig<RecipeStep> {
  recipeID: string;
  recipeStepID: string;

  constructor(recipeID: string, recipeStepID: string, status: number = 200, body?: RecipeStep) {
    super();

    this.recipeID = recipeID;
    this.recipeStepID = recipeStepID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdateRecipeStep = (resCfg: MockUpdateRecipeStepResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps/${resCfg.recipeStepID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
