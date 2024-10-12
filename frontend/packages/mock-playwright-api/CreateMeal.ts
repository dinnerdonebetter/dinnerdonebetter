// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Meal } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockCreateMealResponseConfig extends ResponseConfig<Meal> {
		  

		  constructor(status: number = 201, body?: Meal) {
		    super();

		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockCreateMeal = (resCfg: MockCreateMealResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meals`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};