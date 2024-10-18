// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { UserPermissionsResponse } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCheckPermissionsResponseConfig extends ResponseConfig<UserPermissionsResponse> {
  constructor(status: number = 201, body?: UserPermissionsResponse) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCheckPermissions = (resCfg: MockCheckPermissionsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users/permissions/check`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
