// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { TOTPSecretRefreshResponse } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockRefreshTOTPSecretResponseConfig extends ResponseConfig<TOTPSecretRefreshResponse> {
		  

		  constructor(status: number = 201, body?: TOTPSecretRefreshResponse) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockRefreshTOTPSecret = (resCfg: MockRefreshTOTPSecretResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users/totp_secret/new`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};