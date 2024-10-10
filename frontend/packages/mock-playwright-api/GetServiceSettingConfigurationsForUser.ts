// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ServiceSettingConfiguration, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetServiceSettingConfigurationsForUserResponseConfig extends ResponseConfig<
  QueryFilteredResult<ServiceSettingConfiguration>
> {
  constructor(status: number = 200, body: ServiceSettingConfiguration[] = []) {
    super();

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetServiceSettingConfigurationsForUsers = (
  resCfg: MockGetServiceSettingConfigurationsForUserResponseConfig,
) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/settings/configurations/user`,
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
