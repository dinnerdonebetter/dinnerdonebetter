// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidIngredientGroup } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveValidIngredientGroupResponseConfig extends ResponseConfig<ValidIngredientGroup> {
  validIngredientGroupID: string;

  constructor(validIngredientGroupID: string, status: number = 202, body?: ValidIngredientGroup) {
    super();

    this.validIngredientGroupID = validIngredientGroupID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveValidIngredientGroup = (resCfg: MockArchiveValidIngredientGroupResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_ingredient_groups/${resCfg.validIngredientGroupID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
