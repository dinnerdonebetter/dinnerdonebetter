// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanOptionVote, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetMealPlanOptionVotesResponseConfig extends ResponseConfig<QueryFilteredResult<MealPlanOptionVote>> {
  mealPlanID: string;
  mealPlanEventID: string;
  mealPlanOptionID: string;

  constructor(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    status: number = 200,
    body: MealPlanOptionVote[] = [],
  ) {
    super();

    this.mealPlanID = mealPlanID;
    this.mealPlanEventID = mealPlanEventID;
    this.mealPlanOptionID = mealPlanOptionID;

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetMealPlanOptionVotess = (resCfg: MockGetMealPlanOptionVotesResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/events/${resCfg.mealPlanEventID}/options/${resCfg.mealPlanOptionID}/votes`,
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
