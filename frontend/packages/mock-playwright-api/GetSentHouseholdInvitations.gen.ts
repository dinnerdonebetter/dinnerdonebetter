// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { HouseholdInvitation, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetSentHouseholdInvitationsResponseConfig extends ResponseConfig<
  QueryFilteredResult<HouseholdInvitation>
> {
  constructor(status: number = 200, body: HouseholdInvitation[] = []) {
    super();

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetSentHouseholdInvitationss = (resCfg: MockGetSentHouseholdInvitationsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/household_invitations/sent`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        if (resCfg.body && resCfg.filter) resCfg.body.limit = resCfg.filter.limit;
        if (resCfg.body && resCfg.filter) resCfg.body.page = resCfg.filter.page;

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
