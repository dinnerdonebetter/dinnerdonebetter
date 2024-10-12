// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ServiceSettingConfiguration,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetServiceSettingConfigurationByNameResponseConfig extends ResponseConfig<QueryFilteredResult<ServiceSettingConfiguration>> {
		   serviceSettingConfigurationName: string;
		

		  constructor( serviceSettingConfigurationName: string, status: number = 200, body: ServiceSettingConfiguration[] = []) {
		    super();

		 this.serviceSettingConfigurationName = serviceSettingConfigurationName;
		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetServiceSettingConfigurationByNames = (resCfg: MockGetServiceSettingConfigurationByNameResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/settings/configurations/user/${resCfg.serviceSettingConfigurationName}`,
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