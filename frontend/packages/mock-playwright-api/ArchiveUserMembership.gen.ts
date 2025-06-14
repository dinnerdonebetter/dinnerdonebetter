// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { AccountUserMembership } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveUserMembershipResponseConfig extends ResponseConfig<AccountUserMembership> {
  accountID: string;
  userID: string;

  constructor(accountID: string, userID: string, status: number = 202, body?: AccountUserMembership) {
    super();

    this.accountID = accountID;
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
      `**/api/v1/accounts/${resCfg.accountID}/members/${resCfg.userID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
