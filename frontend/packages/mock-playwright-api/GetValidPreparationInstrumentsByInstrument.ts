// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidPreparationInstrument,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetValidPreparationInstrumentsByInstrumentResponseConfig extends ResponseConfig<QueryFilteredResult<ValidPreparationInstrument>> {
		   validInstrumentID: string;
		

		  constructor( validInstrumentID: string, status: number = 200, body: ValidPreparationInstrument[] = []) {
		    super();

		 this.validInstrumentID = validInstrumentID;
		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetValidPreparationInstrumentsByInstruments = (resCfg: MockGetValidPreparationInstrumentsByInstrumentResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparation_instruments/by_instrument/${resCfg.validInstrumentID}`,
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