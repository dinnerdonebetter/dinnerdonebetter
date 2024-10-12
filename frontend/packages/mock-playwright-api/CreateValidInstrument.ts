// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidInstrument } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateValidInstrumentResponseConfig extends ResponseConfig<ValidInstrument> {
		  

		  constructor(status: number = 201, body?: ValidInstrument) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateValidInstrument = (resCfg: MockCreateValidInstrumentResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_instruments`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};