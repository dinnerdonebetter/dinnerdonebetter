// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { PasswordResetResponse } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdatePasswordResponseConfig extends ResponseConfig<PasswordResetResponse> {
  constructor(status: number = 200, body?: PasswordResetResponse) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdatePassword = (resCfg: MockUpdatePasswordResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users/password/new`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
