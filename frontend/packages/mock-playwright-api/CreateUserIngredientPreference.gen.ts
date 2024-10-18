// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { UserIngredientPreference } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateUserIngredientPreferenceResponseConfig extends ResponseConfig<UserIngredientPreference> {
  constructor(status: number = 201, body?: UserIngredientPreference) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateUserIngredientPreference = (resCfg: MockCreateUserIngredientPreferenceResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/user_ingredient_preferences`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
