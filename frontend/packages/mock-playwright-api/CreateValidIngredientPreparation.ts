// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientPreparation } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockCreateValidIngredientPreparationResponseConfig extends ResponseConfig<ValidIngredientPreparation> {
  constructor(status: number = 201, body?: ValidIngredientPreparation) {
    super();

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockCreateValidIngredientPreparation = (resCfg: MockCreateValidIngredientPreparationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_preparations`,
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
