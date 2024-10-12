// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientStateIngredient } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockArchiveValidIngredientStateIngredientResponseConfig extends ResponseConfig<ValidIngredientStateIngredient> {
		   validIngredientStateIngredientID: string;
		

		  constructor( validIngredientStateIngredientID: string, status: number = 202, body?: ValidIngredientStateIngredient) {
		    super();

		 this.validIngredientStateIngredientID = validIngredientStateIngredientID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockArchiveValidIngredientStateIngredient = (resCfg: MockArchiveValidIngredientStateIngredientResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_state_ingredients/${resCfg.validIngredientStateIngredientID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};