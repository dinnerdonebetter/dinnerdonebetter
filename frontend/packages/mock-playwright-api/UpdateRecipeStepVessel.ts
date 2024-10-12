// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStepVessel } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockUpdateRecipeStepVesselResponseConfig extends ResponseConfig<RecipeStepVessel> {
		   recipeID: string;
		 recipeStepID: string;
		 recipeStepVesselID: string;
		

		  constructor( recipeID: string,  recipeStepID: string,  recipeStepVesselID: string, status: number = 200, body?: RecipeStepVessel) {
		    super();

		 this.recipeID = recipeID;
		 this.recipeStepID = recipeStepID;
		 this.recipeStepVesselID = recipeStepVesselID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockUpdateRecipeStepVessel = (resCfg: MockUpdateRecipeStepVesselResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps/${resCfg.recipeStepID}/vessels/${resCfg.recipeStepVesselID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};