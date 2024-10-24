// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { User } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockVerifyEmailAddressResponseConfig extends ResponseConfig<User> {
  constructor(status: number = 201, body?: User) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockVerifyEmailAddress = (resCfg: MockVerifyEmailAddressResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/users/email_address/verify`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
