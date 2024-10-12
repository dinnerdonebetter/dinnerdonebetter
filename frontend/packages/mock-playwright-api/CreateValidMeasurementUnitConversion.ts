// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidMeasurementUnitConversion } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateValidMeasurementUnitConversionResponseConfig extends ResponseConfig<ValidMeasurementUnitConversion> {
		  

		  constructor(status: number = 201, body?: ValidMeasurementUnitConversion) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateValidMeasurementUnitConversion = (resCfg: MockCreateValidMeasurementUnitConversionResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_measurement_conversions`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};