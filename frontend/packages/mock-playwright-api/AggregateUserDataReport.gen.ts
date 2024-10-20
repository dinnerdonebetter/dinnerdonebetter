// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { UserDataCollectionResponse } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockAggregateUserDataReportResponseConfig extends ResponseConfig<UserDataCollectionResponse> {
  constructor(status: number = 201, body?: UserDataCollectionResponse) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockAggregateUserDataReport = (resCfg: MockAggregateUserDataReportResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/data_privacy/disclose`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
