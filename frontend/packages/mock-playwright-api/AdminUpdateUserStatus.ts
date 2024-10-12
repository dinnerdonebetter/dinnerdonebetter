// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { UserStatusResponse } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockAdminUpdateUserStatusResponseConfig extends ResponseConfig<UserStatusResponse> {
		  

		  constructor(status: number = 201, body?: UserStatusResponse) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockAdminUpdateUserStatus = (resCfg: MockAdminUpdateUserStatusResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/admin/users/status`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};