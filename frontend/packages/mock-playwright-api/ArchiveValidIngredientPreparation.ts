// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientPreparation } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockArchiveValidIngredientPreparationResponseConfig extends ResponseConfig<ValidIngredientPreparation> {
		   validIngredientPreparationID: string;
		

		  constructor( validIngredientPreparationID: string, status: number = 202, body?: ValidIngredientPreparation) {
		    super();

		 this.validIngredientPreparationID = validIngredientPreparationID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockArchiveValidIngredientPreparation = (resCfg: MockArchiveValidIngredientPreparationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_preparations/${resCfg.validIngredientPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};