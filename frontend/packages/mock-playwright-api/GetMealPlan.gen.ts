// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlan } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetMealPlanResponseConfig extends ResponseConfig<MealPlan> {
  mealPlanID: string;

  constructor(mealPlanID: string, status: number = 200, body?: MealPlan) {
    super();

    this.mealPlanID = mealPlanID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetMealPlan = (resCfg: MockGetMealPlanResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
