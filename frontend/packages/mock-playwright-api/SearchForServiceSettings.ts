// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ServiceSetting,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockSearchForServiceSettingsResponseConfig extends ResponseConfig<QueryFilteredResult<ServiceSetting>> {
		   q: string;
		

		  constructor( q: string, status: number = 200, body: ServiceSetting[] = []) {
		    super();

		 this.q = q;
		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockSearchForServiceSettingss = (resCfg: MockSearchForServiceSettingsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/settings/search`,
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