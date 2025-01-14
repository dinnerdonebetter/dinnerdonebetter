// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { TokenResponse } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockAdminLoginForJWTResponseConfig extends ResponseConfig<TokenResponse> {
  constructor(status: number = 201, body?: TokenResponse) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockAdminLoginForJWT = (resCfg: MockAdminLoginForJWTResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/users/login/jwt/admin`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
