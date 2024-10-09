// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { User } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveUserResponseConfig extends ResponseConfig<User> {
  userID: string;

  constructor(userID: string, status: number = 202, body?: User) {
    super();

    this.userID = userID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveUser = (resCfg: MockArchiveUserResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users/${resCfg.userID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
