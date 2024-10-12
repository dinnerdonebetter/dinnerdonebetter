// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { HouseholdInstrumentOwnership } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateHouseholdInstrumentOwnershipResponseConfig extends ResponseConfig<HouseholdInstrumentOwnership> {
		  

		  constructor(status: number = 201, body?: HouseholdInstrumentOwnership) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateHouseholdInstrumentOwnership = (resCfg: MockCreateHouseholdInstrumentOwnershipResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/households/instruments`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};