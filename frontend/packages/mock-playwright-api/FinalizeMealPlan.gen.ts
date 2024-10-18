// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { FinalizeMealPlansResponse } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockFinalizeMealPlanResponseConfig extends ResponseConfig<FinalizeMealPlansResponse> {
  mealPlanID: string;

  constructor(mealPlanID: string, status: number = 201, body?: FinalizeMealPlansResponse) {
    super();

    this.mealPlanID = mealPlanID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockFinalizeMealPlan = (resCfg: MockFinalizeMealPlanResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/finalize`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
