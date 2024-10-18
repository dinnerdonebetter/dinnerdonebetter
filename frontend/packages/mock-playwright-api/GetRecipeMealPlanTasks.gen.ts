// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipePrepTaskStep } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetRecipeMealPlanTasksResponseConfig extends ResponseConfig<RecipePrepTaskStep> {
  recipeID: string;

  constructor(recipeID: string, status: number = 200, body?: RecipePrepTaskStep) {
    super();

    this.recipeID = recipeID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetRecipeMealPlanTasks = (resCfg: MockGetRecipeMealPlanTasksResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/prep_steps`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
