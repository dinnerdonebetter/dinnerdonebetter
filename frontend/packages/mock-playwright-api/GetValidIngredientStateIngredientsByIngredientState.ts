// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientStateIngredient,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetValidIngredientStateIngredientsByIngredientStateResponseConfig extends ResponseConfig<QueryFilteredResult<ValidIngredientStateIngredient>> {
		   validIngredientStateID: string;
		

		  constructor( validIngredientStateID: string, status: number = 200, body: ValidIngredientStateIngredient[] = []) {
		    super();

		 this.validIngredientStateID = validIngredientStateID;
		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetValidIngredientStateIngredientsByIngredientStates = (resCfg: MockGetValidIngredientStateIngredientsByIngredientStateResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/${resCfg.validIngredientStateID}`,
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