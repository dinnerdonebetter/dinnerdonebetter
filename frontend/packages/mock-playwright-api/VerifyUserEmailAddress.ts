// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { User } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockVerifyUserEmailAddressResponseConfig extends ResponseConfig<User> {
		  

		  constructor(status: number = 201, body?: User) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockVerifyUserEmailAddress = (resCfg: MockVerifyUserEmailAddressResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users/email_address_verification`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};