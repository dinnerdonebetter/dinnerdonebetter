// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanGroceryListItem,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetMealPlanGroceryListItemsForMealPlanResponseConfig extends ResponseConfig<QueryFilteredResult<MealPlanGroceryListItem>> {
		   mealPlanID: string;
		

		  constructor( mealPlanID: string, status: number = 200, body: MealPlanGroceryListItem[] = []) {
		    super();

		 this.mealPlanID = mealPlanID;
		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetMealPlanGroceryListItemsForMealPlans = (resCfg: MockGetMealPlanGroceryListItemsForMealPlanResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/grocery_list_items`,
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