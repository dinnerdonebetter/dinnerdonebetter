// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlan, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetMealPlansForHouseholdResponseConfig extends ResponseConfig<QueryFilteredResult<MealPlan>> {
  constructor(status: number = 200, body: MealPlan[] = []) {
    super();

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetMealPlansForHouseholds = (resCfg: MockGetMealPlansForHouseholdResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans`,
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
