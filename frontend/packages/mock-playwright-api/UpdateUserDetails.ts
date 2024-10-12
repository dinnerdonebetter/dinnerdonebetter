// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { User } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockUpdateUserDetailsResponseConfig extends ResponseConfig<User> {
		  

		  constructor(status: number = 200, body?: User) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockUpdateUserDetails = (resCfg: MockUpdateUserDetailsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users/details`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};