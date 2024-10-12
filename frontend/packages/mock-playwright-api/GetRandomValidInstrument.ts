// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidInstrument } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetRandomValidInstrumentResponseConfig extends ResponseConfig<ValidInstrument> {
		  

		  constructor(status: number = 200, body?: ValidInstrument) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetRandomValidInstrument = (resCfg: MockGetRandomValidInstrumentResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_instruments/random`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};