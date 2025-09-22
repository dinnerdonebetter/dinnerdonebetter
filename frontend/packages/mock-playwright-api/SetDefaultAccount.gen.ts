// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Account } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockSetDefaultAccountResponseConfig extends ResponseConfig<Account> {
  accountID: string;

  constructor(accountID: string, status: number = 201, body?: Account) {
    super();

    this.accountID = accountID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockSetDefaultAccount = (resCfg: MockSetDefaultAccountResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/accounts/${resCfg.accountID}/default`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
