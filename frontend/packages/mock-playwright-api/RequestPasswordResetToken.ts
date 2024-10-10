// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { PasswordResetToken } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockRequestPasswordResetTokenResponseConfig extends ResponseConfig<PasswordResetToken> {
  constructor(status: number = 201, body?: PasswordResetToken) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockRequestPasswordResetToken = (resCfg: MockRequestPasswordResetTokenResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/users/password/reset`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
