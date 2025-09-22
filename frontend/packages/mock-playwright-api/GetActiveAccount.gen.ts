// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Account } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetActiveAccountResponseConfig extends ResponseConfig<Account> {
  constructor(status: number = 200, body?: Account) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetActiveAccount = (resCfg: MockGetActiveAccountResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/accounts/current`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
