// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStepCompletionCondition } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdateRecipeStepCompletionConditionResponseConfig extends ResponseConfig<RecipeStepCompletionCondition> {
  recipeID: string;
  recipeStepID: string;
  recipeStepCompletionConditionID: string;

  constructor(
    recipeID: string,
    recipeStepID: string,
    recipeStepCompletionConditionID: string,
    status: number = 200,
    body?: RecipeStepCompletionCondition,
  ) {
    super();

    this.recipeID = recipeID;
    this.recipeStepID = recipeStepID;
    this.recipeStepCompletionConditionID = recipeStepCompletionConditionID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdateRecipeStepCompletionCondition = (
  resCfg: MockUpdateRecipeStepCompletionConditionResponseConfig,
) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps/${resCfg.recipeStepID}/completion_conditions/${resCfg.recipeStepCompletionConditionID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
