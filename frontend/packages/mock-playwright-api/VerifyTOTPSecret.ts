// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { User } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockVerifyTOTPSecretResponseConfig extends ResponseConfig<User> {
		  

		  constructor(status: number = 201, body?: User) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockVerifyTOTPSecret = (resCfg: MockVerifyTOTPSecretResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/users/totp_secret/verify`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};