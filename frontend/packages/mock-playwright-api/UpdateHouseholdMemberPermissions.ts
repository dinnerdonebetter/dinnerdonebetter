// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { UserPermissionsResponse } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockUpdateHouseholdMemberPermissionsResponseConfig extends ResponseConfig<UserPermissionsResponse> {
		   householdID: string;
		 userID: string;
		

		  constructor( householdID: string,  userID: string, status: number = 200, body?: UserPermissionsResponse) {
		    super();

		 this.householdID = householdID;
		 this.userID = userID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockUpdateHouseholdMemberPermissions = (resCfg: MockUpdateHouseholdMemberPermissionsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/households/${resCfg.householdID}/members/${resCfg.userID}/permissions`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PATCH', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};