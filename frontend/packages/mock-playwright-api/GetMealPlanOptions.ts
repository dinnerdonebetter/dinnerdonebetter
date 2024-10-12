// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanOption,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetMealPlanOptionsResponseConfig extends ResponseConfig<QueryFilteredResult<MealPlanOption>> {
		   mealPlanID: string;
		 mealPlanEventID: string;
		

		  constructor( mealPlanID: string,  mealPlanEventID: string, status: number = 200, body: MealPlanOption[] = []) {
		    super();

		 this.mealPlanID = mealPlanID;
		 this.mealPlanEventID = mealPlanEventID;
		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetMealPlanOptionss = (resCfg: MockGetMealPlanOptionsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/events/${resCfg.mealPlanEventID}/options`,
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