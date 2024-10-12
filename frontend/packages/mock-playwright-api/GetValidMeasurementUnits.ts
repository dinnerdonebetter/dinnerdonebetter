// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidMeasurementUnit,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetValidMeasurementUnitsResponseConfig extends ResponseConfig<QueryFilteredResult<ValidMeasurementUnit>> {
		  

		  constructor(status: number = 200, body: ValidMeasurementUnit[] = []) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetValidMeasurementUnitss = (resCfg: MockGetValidMeasurementUnitsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_measurement_units`,
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