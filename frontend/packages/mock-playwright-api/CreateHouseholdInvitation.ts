// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { HouseholdInvitation } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateHouseholdInvitationResponseConfig extends ResponseConfig<HouseholdInvitation> {
		   householdID: string;
		

		  constructor( householdID: string, status: number = 201, body?: HouseholdInvitation) {
		    super();

		 this.householdID = householdID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateHouseholdInvitation = (resCfg: MockCreateHouseholdInvitationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/households/${resCfg.householdID}/invite`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};