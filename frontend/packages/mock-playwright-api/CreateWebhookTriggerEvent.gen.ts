// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { WebhookTriggerEvent } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateWebhookTriggerEventResponseConfig extends ResponseConfig<WebhookTriggerEvent> {
  webhookID: string;

  constructor(webhookID: string, status: number = 201, body?: WebhookTriggerEvent) {
    super();

    this.webhookID = webhookID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateWebhookTriggerEvent = (resCfg: MockCreateWebhookTriggerEventResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/webhooks/${resCfg.webhookID}/trigger_events`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
