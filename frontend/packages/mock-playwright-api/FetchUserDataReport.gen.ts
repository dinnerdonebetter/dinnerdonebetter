// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { UserDataCollection } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockFetchUserDataReportResponseConfig extends ResponseConfig<UserDataCollection> {
  userDataAggregationReportID: string;

  constructor(userDataAggregationReportID: string, status: number = 200, body?: UserDataCollection) {
    super();

    this.userDataAggregationReportID = userDataAggregationReportID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockFetchUserDataReport = (resCfg: MockFetchUserDataReportResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/data_privacy/reports/${resCfg.userDataAggregationReportID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
