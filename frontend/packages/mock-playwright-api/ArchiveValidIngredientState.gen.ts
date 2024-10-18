// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientState } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveValidIngredientStateResponseConfig extends ResponseConfig<ValidIngredientState> {
  validIngredientStateID: string;

  constructor(validIngredientStateID: string, status: number = 202, body?: ValidIngredientState) {
    super();

    this.validIngredientStateID = validIngredientStateID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveValidIngredientState = (resCfg: MockArchiveValidIngredientStateResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_states/${resCfg.validIngredientStateID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
