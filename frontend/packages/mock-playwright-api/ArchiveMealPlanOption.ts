// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanOption } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveMealPlanOptionResponseConfig extends ResponseConfig<MealPlanOption> {
  mealPlanID: string;
  mealPlanEventID: string;
  mealPlanOptionID: string;

  constructor(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    status: number = 202,
    body?: MealPlanOption,
  ) {
    super();

    this.mealPlanID = mealPlanID;
    this.mealPlanEventID = mealPlanEventID;
    this.mealPlanOptionID = mealPlanOptionID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveMealPlanOption = (resCfg: MockArchiveMealPlanOptionResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/events/${resCfg.mealPlanEventID}/options/${resCfg.mealPlanOptionID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
