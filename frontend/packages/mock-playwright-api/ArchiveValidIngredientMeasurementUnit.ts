// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientMeasurementUnit } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockArchiveValidIngredientMeasurementUnitResponseConfig extends ResponseConfig<ValidIngredientMeasurementUnit> {
		   validIngredientMeasurementUnitID: string;
		

		  constructor( validIngredientMeasurementUnitID: string, status: number = 202, body?: ValidIngredientMeasurementUnit) {
		    super();

		 this.validIngredientMeasurementUnitID = validIngredientMeasurementUnitID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockArchiveValidIngredientMeasurementUnit = (resCfg: MockArchiveValidIngredientMeasurementUnitResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_measurement_units/${resCfg.validIngredientMeasurementUnitID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};