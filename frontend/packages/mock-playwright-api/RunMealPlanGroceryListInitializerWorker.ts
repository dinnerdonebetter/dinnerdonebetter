// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { InitializeMealPlanGroceryListResponse } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockRunMealPlanGroceryListInitializerWorkerResponseConfig extends ResponseConfig<InitializeMealPlanGroceryListResponse> {
		  

		  constructor(status: number = 201, body?: InitializeMealPlanGroceryListResponse) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockRunMealPlanGroceryListInitializerWorker = (resCfg: MockRunMealPlanGroceryListInitializerWorkerResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/workers/meal_plan_grocery_list_init`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};