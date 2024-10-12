// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ServiceSetting } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateServiceSettingResponseConfig extends ResponseConfig<ServiceSetting> {
		  

		  constructor(status: number = 201, body?: ServiceSetting) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateServiceSetting = (resCfg: MockCreateServiceSettingResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/settings`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};