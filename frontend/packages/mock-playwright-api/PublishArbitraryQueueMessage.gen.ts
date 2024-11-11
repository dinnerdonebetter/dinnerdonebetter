// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ArbitraryQueueMessageResponse } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockPublishArbitraryQueueMessageResponseConfig extends ResponseConfig<ArbitraryQueueMessageResponse> {
  constructor(status: number = 201, body?: ArbitraryQueueMessageResponse) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockPublishArbitraryQueueMessage = (resCfg: MockPublishArbitraryQueueMessageResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/admin/queues/test`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
