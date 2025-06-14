// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Account } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdateAccountResponseConfig extends ResponseConfig<Account> {
  accountID: string;

  constructor(accountID: string, status: number = 200, body?: Account) {
    super();

    this.accountID = accountID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdateAccount = (resCfg: MockUpdateAccountResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/accounts/${resCfg.accountID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
