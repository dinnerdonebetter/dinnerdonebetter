// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlan } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateMealPlanResponseConfig extends ResponseConfig<MealPlan> {
		  

		  constructor(status: number = 201, body?: MealPlan) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateMealPlan = (resCfg: MockCreateMealPlanResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};