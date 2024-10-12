// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidInstrument } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockUpdateValidInstrumentResponseConfig extends ResponseConfig<ValidInstrument> {
		   validInstrumentID: string;
		

		  constructor( validInstrumentID: string, status: number = 200, body?: ValidInstrument) {
		    super();

		 this.validInstrumentID = validInstrumentID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockUpdateValidInstrument = (resCfg: MockUpdateValidInstrumentResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_instruments/${resCfg.validInstrumentID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};