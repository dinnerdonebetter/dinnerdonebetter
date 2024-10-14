// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanOptionVote } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveMealPlanOptionVoteResponseConfig extends ResponseConfig<MealPlanOptionVote> {
  mealPlanID: string;
  mealPlanEventID: string;
  mealPlanOptionID: string;
  mealPlanOptionVoteID: string;

  constructor(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    mealPlanOptionVoteID: string,
    status: number = 202,
    body?: MealPlanOptionVote,
  ) {
    super();

    this.mealPlanID = mealPlanID;
    this.mealPlanEventID = mealPlanEventID;
    this.mealPlanOptionID = mealPlanOptionID;
    this.mealPlanOptionVoteID = mealPlanOptionVoteID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveMealPlanOptionVote = (resCfg: MockArchiveMealPlanOptionVoteResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/events/${resCfg.mealPlanEventID}/options/${resCfg.mealPlanOptionID}/votes/${resCfg.mealPlanOptionVoteID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
