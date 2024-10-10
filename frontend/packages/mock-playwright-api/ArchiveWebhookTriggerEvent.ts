// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { WebhookTriggerEvent } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveWebhookTriggerEventResponseConfig extends ResponseConfig<WebhookTriggerEvent> {
  webhookID: string;
  webhookTriggerEventID: string;

  constructor(webhookID: string, webhookTriggerEventID: string, status: number = 202, body?: WebhookTriggerEvent) {
    super();

    this.webhookID = webhookID;
    this.webhookTriggerEventID = webhookTriggerEventID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveWebhookTriggerEvent = (resCfg: MockArchiveWebhookTriggerEventResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/webhooks/${resCfg.webhookID}/trigger_events/${resCfg.webhookTriggerEventID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
