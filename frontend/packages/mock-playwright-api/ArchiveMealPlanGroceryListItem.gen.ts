// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { MealPlanGroceryListItem } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveMealPlanGroceryListItemResponseConfig extends ResponseConfig<MealPlanGroceryListItem> {
  mealPlanID: string;
  mealPlanGroceryListItemID: string;

  constructor(
    mealPlanID: string,
    mealPlanGroceryListItemID: string,
    status: number = 202,
    body?: MealPlanGroceryListItem,
  ) {
    super();

    this.mealPlanID = mealPlanID;
    this.mealPlanGroceryListItemID = mealPlanGroceryListItemID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveMealPlanGroceryListItem = (resCfg: MockArchiveMealPlanGroceryListItemResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meal_plans/${resCfg.mealPlanID}/grocery_list_items/${resCfg.mealPlanGroceryListItemID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
