// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStepIngredient } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockArchiveRecipeStepIngredientResponseConfig extends ResponseConfig<RecipeStepIngredient> {
		   recipeID: string;
		 recipeStepID: string;
		 recipeStepIngredientID: string;
		

		  constructor( recipeID: string,  recipeStepID: string,  recipeStepIngredientID: string, status: number = 202, body?: RecipeStepIngredient) {
		    super();

		 this.recipeID = recipeID;
		 this.recipeStepID = recipeStepID;
		 this.recipeStepIngredientID = recipeStepIngredientID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockArchiveRecipeStepIngredient = (resCfg: MockArchiveRecipeStepIngredientResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps/${resCfg.recipeStepID}/ingredients/${resCfg.recipeStepIngredientID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};