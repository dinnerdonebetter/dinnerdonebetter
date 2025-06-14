// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { AccountInvitation, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetReceivedAccountInvitationsResponseConfig extends ResponseConfig<
  QueryFilteredResult<AccountInvitation>
> {
  constructor(status: number = 200, body: AccountInvitation[] = []) {
    super();

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetReceivedAccountInvitationss = (resCfg: MockGetReceivedAccountInvitationsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/account_invitations/received`,
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
