// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { OAuth2Client,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetOAuth2ClientsResponseConfig extends ResponseConfig<QueryFilteredResult<OAuth2Client>> {
		  

		  constructor(status: number = 200, body: OAuth2Client[] = []) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetOAuth2Clientss = (resCfg: MockGetOAuth2ClientsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/oauth2_clients`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		
        if (resCfg.body && resCfg.filter) resCfg.body.limit = resCfg.filter.limit;
        if (resCfg.body && resCfg.filter) resCfg.body.page = resCfg.filter.page;
		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};