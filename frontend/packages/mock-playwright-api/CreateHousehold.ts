// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Household } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateHouseholdResponseConfig extends ResponseConfig<Household> {
		  

		  constructor(status: number = 201, body?: Household) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateHousehold = (resCfg: MockCreateHouseholdResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/households`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};