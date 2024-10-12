// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { User } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetSelfResponseConfig extends ResponseConfig<User> {
		  

		  constructor(status: number = 200, body?: User) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetSelf = (resCfg: MockGetSelfResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users/self`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};