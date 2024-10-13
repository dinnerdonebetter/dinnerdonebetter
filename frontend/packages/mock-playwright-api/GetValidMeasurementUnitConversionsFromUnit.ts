// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidMeasurementUnitConversion,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetValidMeasurementUnitConversionsFromUnitResponseConfig extends ResponseConfig<QueryFilteredResult<ValidMeasurementUnitConversion>> {
		   validMeasurementUnitID: string;
		

		  constructor( validMeasurementUnitID: string, status: number = 200, body: ValidMeasurementUnitConversion[] = []) {
		    super();

		 this.validMeasurementUnitID = validMeasurementUnitID;
		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetValidMeasurementUnitConversionsFromUnits = (resCfg: MockGetValidMeasurementUnitConversionsFromUnitResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_measurement_conversions/from_unit/${resCfg.validMeasurementUnitID}`,
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