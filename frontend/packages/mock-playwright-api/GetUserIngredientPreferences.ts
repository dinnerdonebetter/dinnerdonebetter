// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { UserIngredientPreference,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetUserIngredientPreferencesResponseConfig extends ResponseConfig<QueryFilteredResult<UserIngredientPreference>> {
		  

		  constructor(status: number = 200, body: UserIngredientPreference[] = []) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetUserIngredientPreferencess = (resCfg: MockGetUserIngredientPreferencesResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/user_ingredient_preferences`,
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