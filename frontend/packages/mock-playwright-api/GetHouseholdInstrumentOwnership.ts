// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { HouseholdInstrumentOwnership } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetHouseholdInstrumentOwnershipResponseConfig extends ResponseConfig<HouseholdInstrumentOwnership> {
		   householdInstrumentOwnershipID: string;
		

		  constructor( householdInstrumentOwnershipID: string, status: number = 200, body?: HouseholdInstrumentOwnership) {
		    super();

		 this.householdInstrumentOwnershipID = householdInstrumentOwnershipID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetHouseholdInstrumentOwnership = (resCfg: MockGetHouseholdInstrumentOwnershipResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/households/instruments/${resCfg.householdInstrumentOwnershipID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};