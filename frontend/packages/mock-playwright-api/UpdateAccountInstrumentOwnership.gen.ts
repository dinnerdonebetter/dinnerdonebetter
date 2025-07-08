// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { AccountInstrumentOwnership } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdateAccountInstrumentOwnershipResponseConfig extends ResponseConfig<AccountInstrumentOwnership> {
  accountInstrumentOwnershipID: string;

  constructor(accountInstrumentOwnershipID: string, status: number = 200, body?: AccountInstrumentOwnership) {
    super();

    this.accountInstrumentOwnershipID = accountInstrumentOwnershipID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdateAccountInstrumentOwnership = (
  resCfg: MockUpdateAccountInstrumentOwnershipResponseConfig,
) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/accounts/instruments/${resCfg.accountInstrumentOwnershipID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
