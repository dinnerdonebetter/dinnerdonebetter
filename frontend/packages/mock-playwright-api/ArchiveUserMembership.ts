// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { HouseholdUserMembership } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveUserMembershipResponseConfig extends ResponseConfig<HouseholdUserMembership> {
  householdID: string;
  userID: string;

  constructor(householdID: string, userID: string, status: number = 202, body?: HouseholdUserMembership) {
    super();

    this.householdID = householdID;
    this.userID = userID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveUserMembership = (resCfg: MockArchiveUserMembershipResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/households/${resCfg.householdID}/members/${resCfg.userID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
