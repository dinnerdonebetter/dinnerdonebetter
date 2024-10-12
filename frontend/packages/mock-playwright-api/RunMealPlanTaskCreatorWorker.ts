// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { CreateMealPlanTasksResponse } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockRunMealPlanTaskCreatorWorkerResponseConfig extends ResponseConfig<CreateMealPlanTasksResponse> {
		  

		  constructor(status: number = 201, body?: CreateMealPlanTasksResponse) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockRunMealPlanTaskCreatorWorker = (resCfg: MockRunMealPlanTaskCreatorWorkerResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/workers/meal_plan_tasks`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};