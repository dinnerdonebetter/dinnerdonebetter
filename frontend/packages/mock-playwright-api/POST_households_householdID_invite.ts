// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { HouseholdInvitation } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockPOST_households_householdID_inviteResponseConfig extends ResponseConfig<HouseholdInvitation> {
  householdID: string;

  constructor(householdID: string, status: number = 201, body?: HouseholdInvitation) {
    super();

    this.householdID = householdID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockPOST_households_householdID_invite = (
  resCfg: MockPOST_households_householdID_inviteResponseConfig,
) => {
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
