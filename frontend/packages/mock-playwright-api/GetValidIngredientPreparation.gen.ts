// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientPreparation } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockGetValidIngredientPreparationResponseConfig extends ResponseConfig<ValidIngredientPreparation> {
  validIngredientPreparationID: string;

  constructor(validIngredientPreparationID: string, status: number = 200, body?: ValidIngredientPreparation) {
    super();

    this.validIngredientPreparationID = validIngredientPreparationID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockGetValidIngredientPreparation = (resCfg: MockGetValidIngredientPreparationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_preparations/${resCfg.validIngredientPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
