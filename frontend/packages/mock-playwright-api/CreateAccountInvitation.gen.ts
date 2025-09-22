// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { AccountInvitation } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateAccountInvitationResponseConfig extends ResponseConfig<AccountInvitation> {
  accountID: string;

  constructor(accountID: string, status: number = 201, body?: AccountInvitation) {
    super();

    this.accountID = accountID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateAccountInvitation = (resCfg: MockCreateAccountInvitationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/accounts/${resCfg.accountID}/invite`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
