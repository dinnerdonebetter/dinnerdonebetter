// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { AccountInstrumentOwnership } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateAccountInstrumentOwnershipResponseConfig extends ResponseConfig<AccountInstrumentOwnership> {
  constructor(status: number = 201, body?: AccountInstrumentOwnership) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateAccountInstrumentOwnership = (
  resCfg: MockCreateAccountInstrumentOwnershipResponseConfig,
) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/accounts/instruments`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
