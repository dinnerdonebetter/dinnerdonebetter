// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanGroceryListItem } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateMealPlanGroceryListItemResponseConfig extends ResponseConfig<MealPlanGroceryListItem> {
  mealPlanID: string;

  constructor(mealPlanID: string, status: number = 201, body?: MealPlanGroceryListItem) {
    super();

    this.mealPlanID = mealPlanID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateMealPlanGroceryListItem = (resCfg: MockCreateMealPlanGroceryListItemResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/grocery_list_items`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
