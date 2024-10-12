// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipeStepCompletionCondition } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockArchiveRecipeStepCompletionConditionResponseConfig extends ResponseConfig<RecipeStepCompletionCondition> {
		   recipeID: string;
		 recipeStepID: string;
		 recipeStepCompletionConditionID: string;
		

		  constructor( recipeID: string,  recipeStepID: string,  recipeStepCompletionConditionID: string, status: number = 202, body?: RecipeStepCompletionCondition) {
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

export const mockArchiveRecipeStepCompletionCondition = (resCfg: MockArchiveRecipeStepCompletionConditionResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/steps/${resCfg.recipeStepID}/completion_conditions/${resCfg.recipeStepCompletionConditionID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};