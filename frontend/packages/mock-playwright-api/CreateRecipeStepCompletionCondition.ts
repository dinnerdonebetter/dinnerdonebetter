// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStepCompletionCondition } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateRecipeStepCompletionConditionResponseConfig extends ResponseConfig<RecipeStepCompletionCondition> {
  recipeID: string;
  recipeStepID: string;

  constructor(recipeID: string, recipeStepID: string, status: number = 201, body?: RecipeStepCompletionCondition) {
    super();

    this.recipeID = recipeID;
    this.recipeStepID = recipeStepID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateRecipeStepCompletionCondition = (
  resCfg: MockCreateRecipeStepCompletionConditionResponseConfig,
) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps/${resCfg.recipeStepID}/completion_conditions`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
