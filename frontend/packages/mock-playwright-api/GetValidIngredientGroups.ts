// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientGroup,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetValidIngredientGroupsResponseConfig extends ResponseConfig<QueryFilteredResult<ValidIngredientGroup>> {
		  

		  constructor(status: number = 200, body: ValidIngredientGroup[] = []) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetValidIngredientGroupss = (resCfg: MockGetValidIngredientGroupsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_groups`,
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