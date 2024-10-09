// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Webhook } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateWebhookResponseConfig extends ResponseConfig<Webhook> {
  constructor(status: number = 201, body?: Webhook) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateWebhook = (resCfg: MockCreateWebhookResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/webhooks`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
