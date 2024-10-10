// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredient } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockUpdateValidIngredientResponseConfig extends ResponseConfig<ValidIngredient> {
  validIngredientID: string;

  constructor(validIngredientID: string, status: number = 200, body?: ValidIngredient) {
    super();

    this.validIngredientID = validIngredientID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockUpdateValidIngredient = (resCfg: MockUpdateValidIngredientResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredients/${resCfg.validIngredientID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
