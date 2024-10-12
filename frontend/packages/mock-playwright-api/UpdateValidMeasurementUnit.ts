// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidMeasurementUnit } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockUpdateValidMeasurementUnitResponseConfig extends ResponseConfig<ValidMeasurementUnit> {
		   validMeasurementUnitID: string;
		

		  constructor( validMeasurementUnitID: string, status: number = 200, body?: ValidMeasurementUnit) {
		    super();

		 this.validMeasurementUnitID = validMeasurementUnitID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockUpdateValidMeasurementUnit = (resCfg: MockUpdateValidMeasurementUnitResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_measurement_units/${resCfg.validMeasurementUnitID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};