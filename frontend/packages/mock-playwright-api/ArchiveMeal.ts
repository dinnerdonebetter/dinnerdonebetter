// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { Meal } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveMealResponseConfig extends ResponseConfig<Meal> {
  mealID: string;

  constructor(mealID: string, status: number = 202, body?: Meal) {
    super();

    this.mealID = mealID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveMeal = (resCfg: MockArchiveMealResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/meals/${resCfg.mealID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
