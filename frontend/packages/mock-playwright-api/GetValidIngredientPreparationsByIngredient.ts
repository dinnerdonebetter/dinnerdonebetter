// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientPreparation,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetValidIngredientPreparationsByIngredientResponseConfig extends ResponseConfig<QueryFilteredResult<ValidIngredientPreparation>> {
		   validIngredientID: string;
		

		  constructor( validIngredientID: string, status: number = 200, body: ValidIngredientPreparation[] = []) {
		    super();

		 this.validIngredientID = validIngredientID;
		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetValidIngredientPreparationsByIngredients = (resCfg: MockGetValidIngredientPreparationsByIngredientResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_preparations/by_ingredient/${resCfg.validIngredientID}`,
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