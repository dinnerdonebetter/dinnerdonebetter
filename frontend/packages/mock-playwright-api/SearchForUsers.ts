// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { User,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockSearchForUsersResponseConfig extends ResponseConfig<QueryFilteredResult<User>> {
		   q: string;
		

		  constructor( q: string, status: number = 200, body: User[] = []) {
		    super();

		 this.q = q;
		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockSearchForUserss = (resCfg: MockSearchForUsersResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users/search`,
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