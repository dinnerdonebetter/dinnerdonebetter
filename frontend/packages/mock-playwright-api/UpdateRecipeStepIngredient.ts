// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStepIngredient } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockUpdateRecipeStepIngredientResponseConfig extends ResponseConfig<RecipeStepIngredient> {
		   recipeID: string;
		 recipeStepID: string;
		 recipeStepIngredientID: string;
		

		  constructor( recipeID: string,  recipeStepID: string,  recipeStepIngredientID: string, status: number = 200, body?: RecipeStepIngredient) {
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

export const mockUpdateRecipeStepIngredient = (resCfg: MockUpdateRecipeStepIngredientResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps/${resCfg.recipeStepID}/ingredients/${resCfg.recipeStepIngredientID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};