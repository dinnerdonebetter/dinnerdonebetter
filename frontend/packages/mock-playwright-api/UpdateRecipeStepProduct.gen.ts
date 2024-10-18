// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStepProduct } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdateRecipeStepProductResponseConfig extends ResponseConfig<RecipeStepProduct> {
  recipeID: string;
  recipeStepID: string;
  recipeStepProductID: string;

  constructor(
    recipeID: string,
    recipeStepID: string,
    recipeStepProductID: string,
    status: number = 200,
    body?: RecipeStepProduct,
  ) {
    super();

    this.recipeID = recipeID;
    this.recipeStepID = recipeStepID;
    this.recipeStepProductID = recipeStepProductID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdateRecipeStepProduct = (resCfg: MockUpdateRecipeStepProductResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps/${resCfg.recipeStepID}/products/${resCfg.recipeStepProductID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
