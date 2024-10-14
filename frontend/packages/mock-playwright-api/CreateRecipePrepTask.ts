// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipePrepTask } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateRecipePrepTaskResponseConfig extends ResponseConfig<RecipePrepTask> {
  recipeID: string;

  constructor(recipeID: string, status: number = 201, body?: RecipePrepTask) {
    super();

    this.recipeID = recipeID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateRecipePrepTask = (resCfg: MockCreateRecipePrepTaskResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/prep_tasks`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
