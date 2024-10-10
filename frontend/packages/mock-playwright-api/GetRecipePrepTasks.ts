// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipePrepTask, QueryFilteredResult } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetRecipePrepTasksResponseConfig extends ResponseConfig<QueryFilteredResult<RecipePrepTask>> {
  recipeID: string;

  constructor(recipeID: string, status: number = 200, body: RecipePrepTask[] = []) {
    super();

    this.recipeID = recipeID;

    this.status = status;
    if (this.body) {
      this.body.data = body;
    }
  }
}

export const mockGetRecipePrepTaskss = (resCfg: MockGetRecipePrepTasksResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/prep_tasks`,
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
