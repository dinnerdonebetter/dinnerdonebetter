// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Webhook } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetWebhookResponseConfig extends ResponseConfig<Webhook> {
		   webhookID: string;
		

		  constructor( webhookID: string, status: number = 200, body?: Webhook) {
		    super();

		 this.webhookID = webhookID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetWebhook = (resCfg: MockGetWebhookResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/webhooks/${resCfg.webhookID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};