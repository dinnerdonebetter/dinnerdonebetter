// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { AccountInvitation } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetAccountInvitationByIDResponseConfig extends ResponseConfig<AccountInvitation> {
  accountID: string;
  accountInvitationID: string;

  constructor(accountID: string, accountInvitationID: string, status: number = 200, body?: AccountInvitation) {
    super();

    this.accountID = accountID;
    this.accountInvitationID = accountInvitationID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetAccountInvitationByID = (resCfg: MockGetAccountInvitationByIDResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/accounts/${resCfg.accountID}/invitations/${resCfg.accountInvitationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
