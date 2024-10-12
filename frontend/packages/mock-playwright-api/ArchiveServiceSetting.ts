// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ServiceSetting } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockArchiveServiceSettingResponseConfig extends ResponseConfig<ServiceSetting> {
		   serviceSettingID: string;
		

		  constructor( serviceSettingID: string, status: number = 202, body?: ServiceSetting) {
		    super();

		 this.serviceSettingID = serviceSettingID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockArchiveServiceSetting = (resCfg: MockArchiveServiceSettingResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/settings/${resCfg.serviceSettingID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};