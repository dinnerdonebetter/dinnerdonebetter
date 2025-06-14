// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { UserPermissionsResponse } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdateAccountMemberPermissionsResponseConfig extends ResponseConfig<UserPermissionsResponse> {
  accountID: string;
  userID: string;

  constructor(accountID: string, userID: string, status: number = 200, body?: UserPermissionsResponse) {
    super();

    this.accountID = accountID;
    this.userID = userID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdateAccountMemberPermissions = (resCfg: MockUpdateAccountMemberPermissionsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/accounts/${resCfg.accountID}/members/${resCfg.userID}/permissions`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PATCH', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
