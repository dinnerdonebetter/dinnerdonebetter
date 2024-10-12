// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { FinalizeMealPlansResponse } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockRunFinalizeMealPlanWorkerResponseConfig extends ResponseConfig<FinalizeMealPlansResponse> {
		  

		  constructor(status: number = 201, body?: FinalizeMealPlansResponse) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockRunFinalizeMealPlanWorker = (resCfg: MockRunFinalizeMealPlanWorkerResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/workers/finalize_meal_plans`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};