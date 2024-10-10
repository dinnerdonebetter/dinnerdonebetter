// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredient } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveValidIngredientResponseConfig extends ResponseConfig<ValidIngredient> {
  validIngredientID: string;

  constructor(validIngredientID: string, status: number = 202, body?: ValidIngredient) {
    super();

    this.validIngredientID = validIngredientID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveValidIngredient = (resCfg: MockArchiveValidIngredientResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredients/${resCfg.validIngredientID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
