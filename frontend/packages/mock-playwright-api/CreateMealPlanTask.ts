// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanTask } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateMealPlanTaskResponseConfig extends ResponseConfig<MealPlanTask> {
		   mealPlanID: string;
		

		  constructor( mealPlanID: string, status: number = 201, body?: MealPlanTask) {
		    super();

		 this.mealPlanID = mealPlanID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateMealPlanTask = (resCfg: MockCreateMealPlanTaskResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/tasks`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};