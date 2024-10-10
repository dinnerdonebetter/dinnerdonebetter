// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ServiceSettingConfiguration } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdateServiceSettingConfigurationResponseConfig extends ResponseConfig<ServiceSettingConfiguration> {
  serviceSettingConfigurationID: string;

  constructor(serviceSettingConfigurationID: string, status: number = 200, body?: ServiceSettingConfiguration) {
    super();

    this.serviceSettingConfigurationID = serviceSettingConfigurationID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdateServiceSettingConfiguration = (resCfg: MockUpdateServiceSettingConfigurationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/settings/configurations/${resCfg.serviceSettingConfigurationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
