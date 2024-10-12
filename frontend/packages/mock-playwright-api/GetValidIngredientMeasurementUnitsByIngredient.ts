// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientMeasurementUnit,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetValidIngredientMeasurementUnitsByIngredientResponseConfig extends ResponseConfig<QueryFilteredResult<ValidIngredientMeasurementUnit>> {
		   validIngredientID: string;
		

		  constructor( validIngredientID: string, status: number = 200, body: ValidIngredientMeasurementUnit[] = []) {
		    super();

		 this.validIngredientID = validIngredientID;
		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetValidIngredientMeasurementUnitsByIngredients = (resCfg: MockGetValidIngredientMeasurementUnitsByIngredientResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_measurement_units/by_ingredient/${resCfg.validIngredientID}`,
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