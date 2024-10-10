// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanEvent } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateMealPlanEventResponseConfig extends ResponseConfig<MealPlanEvent> {
  mealPlanID: string;

  constructor(mealPlanID: string, status: number = 201, body?: MealPlanEvent) {
    super();

    this.mealPlanID = mealPlanID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateMealPlanEvent = (resCfg: MockCreateMealPlanEventResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/events`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
