// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanOptionVote } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateMealPlanOptionVoteResponseConfig extends ResponseConfig<MealPlanOptionVote> {
		   mealPlanID: string;
		 mealPlanEventID: string;
		

		  constructor( mealPlanID: string,  mealPlanEventID: string, status: number = 201, body?: MealPlanOptionVote) {
		    super();

		 this.mealPlanID = mealPlanID;
		 this.mealPlanEventID = mealPlanEventID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateMealPlanOptionVote = (resCfg: MockCreateMealPlanOptionVoteResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/events/${resCfg.mealPlanEventID}/vote`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};