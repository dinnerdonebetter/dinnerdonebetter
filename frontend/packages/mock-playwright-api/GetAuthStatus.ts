// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { UserStatusResponse } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetAuthStatusResponseConfig extends ResponseConfig<UserStatusResponse> {
  constructor(status: number = 200, body?: UserStatusResponse) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetAuthStatus = (resCfg: MockGetAuthStatusResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/auth/status`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
