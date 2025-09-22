// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Account } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateAccountResponseConfig extends ResponseConfig<Account> {
  constructor(status: number = 201, body?: Account) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateAccount = (resCfg: MockCreateAccountResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/accounts`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
