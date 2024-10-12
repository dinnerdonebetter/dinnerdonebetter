// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Household } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockTransferHouseholdOwnershipResponseConfig extends ResponseConfig<Household> {
		   householdID: string;
		

		  constructor( householdID: string, status: number = 201, body?: Household) {
		    super();

		 this.householdID = householdID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockTransferHouseholdOwnership = (resCfg: MockTransferHouseholdOwnershipResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/households/${resCfg.householdID}/transfer`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};