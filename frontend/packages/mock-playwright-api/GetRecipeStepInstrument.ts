// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStepInstrument } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetRecipeStepInstrumentResponseConfig extends ResponseConfig<RecipeStepInstrument> {
		   recipeID: string;
		 recipeStepID: string;
		 recipeStepInstrumentID: string;
		

		  constructor( recipeID: string,  recipeStepID: string,  recipeStepInstrumentID: string, status: number = 200, body?: RecipeStepInstrument) {
		    super();

		 this.recipeID = recipeID;
		 this.recipeStepID = recipeStepID;
		 this.recipeStepInstrumentID = recipeStepInstrumentID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetRecipeStepInstrument = (resCfg: MockGetRecipeStepInstrumentResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps/${resCfg.recipeStepID}/instruments/${resCfg.recipeStepInstrumentID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};