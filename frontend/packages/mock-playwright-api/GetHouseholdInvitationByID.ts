// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { HouseholdInvitation } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetHouseholdInvitationByIDResponseConfig extends ResponseConfig<HouseholdInvitation> {
		   householdID: string;
		 householdInvitationID: string;
		

		  constructor( householdID: string,  householdInvitationID: string, status: number = 200, body?: HouseholdInvitation) {
		    super();

		 this.householdID = householdID;
		 this.householdInvitationID = householdInvitationID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetHouseholdInvitationByID = (resCfg: MockGetHouseholdInvitationByIDResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/households/${resCfg.householdID}/invitations/${resCfg.householdInvitationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};