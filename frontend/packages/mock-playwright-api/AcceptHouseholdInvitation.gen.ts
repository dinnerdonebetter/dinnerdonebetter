// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { HouseholdInvitation } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockAcceptHouseholdInvitationResponseConfig extends ResponseConfig<HouseholdInvitation> {
  householdInvitationID: string;

  constructor(householdInvitationID: string, status: number = 200, body?: HouseholdInvitation) {
    super();

    this.householdInvitationID = householdInvitationID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockAcceptHouseholdInvitation = (resCfg: MockAcceptHouseholdInvitationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/household_invitations/${resCfg.householdInvitationID}/accept`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
