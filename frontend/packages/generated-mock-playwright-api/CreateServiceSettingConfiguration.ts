// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ServiceSettingConfiguration } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateServiceSettingConfigurationResponseConfig extends ResponseConfig<ServiceSettingConfiguration> {
  constructor(status: number = 201, body?: ServiceSettingConfiguration) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateServiceSettingConfiguration = (resCfg: MockCreateServiceSettingConfigurationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/settings/configurations`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
