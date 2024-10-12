// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { OAuth2Client } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetOAuth2ClientResponseConfig extends ResponseConfig<OAuth2Client> {
		   oauth2ClientID: string;
		

		  constructor( oauth2ClientID: string, status: number = 200, body?: OAuth2Client) {
		    super();

		 this.oauth2ClientID = oauth2ClientID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetOAuth2Client = (resCfg: MockGetOAuth2ClientResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/oauth2_clients/${resCfg.oauth2ClientID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};