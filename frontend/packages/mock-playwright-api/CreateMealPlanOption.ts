// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanOption } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateMealPlanOptionResponseConfig extends ResponseConfig<MealPlanOption> {
  mealPlanID: string;
  mealPlanEventID: string;

  constructor(mealPlanID: string, mealPlanEventID: string, status: number = 201, body?: MealPlanOption) {
    super();

    this.mealPlanID = mealPlanID;
    this.mealPlanEventID = mealPlanEventID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateMealPlanOption = (resCfg: MockCreateMealPlanOptionResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/events/${resCfg.mealPlanEventID}/options`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
