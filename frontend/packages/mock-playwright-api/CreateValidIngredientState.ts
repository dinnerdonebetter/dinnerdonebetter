// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientState } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateValidIngredientStateResponseConfig extends ResponseConfig<ValidIngredientState> {
		  

		  constructor(status: number = 201, body?: ValidIngredientState) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateValidIngredientState = (resCfg: MockCreateValidIngredientStateResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_states`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};