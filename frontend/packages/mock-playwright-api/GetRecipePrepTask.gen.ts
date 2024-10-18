// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { RecipePrepTask } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetRecipePrepTaskResponseConfig extends ResponseConfig<RecipePrepTask> {
  recipeID: string;
  recipePrepTaskID: string;

  constructor(recipeID: string, recipePrepTaskID: string, status: number = 200, body?: RecipePrepTask) {
    super();

    this.recipeID = recipeID;
    this.recipePrepTaskID = recipePrepTaskID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetRecipePrepTask = (resCfg: MockGetRecipePrepTaskResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/recipes/${resCfg.recipeID}/prep_tasks/${resCfg.recipePrepTaskID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
