// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanEvent } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetMealPlanEventResponseConfig extends ResponseConfig<MealPlanEvent> {
		   mealPlanID: string;
		 mealPlanEventID: string;
		

		  constructor( mealPlanID: string,  mealPlanEventID: string, status: number = 200, body?: MealPlanEvent) {
		    super();

		 this.mealPlanID = mealPlanID;
		 this.mealPlanEventID = mealPlanEventID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetMealPlanEvent = (resCfg: MockGetMealPlanEventResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/events/${resCfg.mealPlanEventID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};