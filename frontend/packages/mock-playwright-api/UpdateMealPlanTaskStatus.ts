// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanTask } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockUpdateMealPlanTaskStatusResponseConfig extends ResponseConfig<MealPlanTask> {
		   mealPlanID: string;
		 mealPlanTaskID: string;
		

		  constructor( mealPlanID: string,  mealPlanTaskID: string, status: number = 200, body?: MealPlanTask) {
		    super();

		 this.mealPlanID = mealPlanID;
		 this.mealPlanTaskID = mealPlanTaskID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockUpdateMealPlanTaskStatus = (resCfg: MockUpdateMealPlanTaskStatusResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/tasks/${resCfg.mealPlanTaskID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PATCH', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};