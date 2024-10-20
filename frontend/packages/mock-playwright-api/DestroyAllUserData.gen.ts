// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { DataDeletionResponse } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockDestroyAllUserDataResponseConfig extends ResponseConfig<DataDeletionResponse> {
  constructor(status: number = 202, body?: DataDeletionResponse) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockDestroyAllUserData = (resCfg: MockDestroyAllUserDataResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/data_privacy/destroy`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
