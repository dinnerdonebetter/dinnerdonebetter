// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidPreparation } from '@dinnerdonebetter/models';

import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockArchiveValidPreparationResponseConfig extends ResponseConfig<ValidPreparation> {
  validPreparationID: string;

  constructor(validPreparationID: string, status: number = 202, body?: ValidPreparation) {
    super();

    this.validPreparationID = validPreparationID;

    this.status = status;
    if (this.body) {
      this.body = body;
    }
  }
}

export const mockArchiveValidPreparation = (resCfg: MockArchiveValidPreparationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparations/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
