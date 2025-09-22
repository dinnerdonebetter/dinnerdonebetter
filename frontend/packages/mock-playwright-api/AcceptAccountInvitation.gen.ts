// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { AccountInvitation } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockAcceptAccountInvitationResponseConfig extends ResponseConfig<AccountInvitation> {
  accountInvitationID: string;

  constructor(accountInvitationID: string, status: number = 200, body?: AccountInvitation) {
    super();

    this.accountInvitationID = accountInvitationID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockAcceptAccountInvitation = (resCfg: MockAcceptAccountInvitationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/account_invitations/${resCfg.accountInvitationID}/accept`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
