// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { User } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetUserResponseConfig extends ResponseConfig<User> {
  userID: string;

  constructor(userID: string, status: number = 200, body?: User) {
    super();

    this.userID = userID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetUser = (resCfg: MockGetUserResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users/${resCfg.userID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
