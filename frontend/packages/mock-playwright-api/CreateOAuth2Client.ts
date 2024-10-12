// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { OAuth2ClientCreationResponse } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateOAuth2ClientResponseConfig extends ResponseConfig<OAuth2ClientCreationResponse> {
		  

		  constructor(status: number = 201, body?: OAuth2ClientCreationResponse) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateOAuth2Client = (resCfg: MockCreateOAuth2ClientResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/oauth2_clients`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};