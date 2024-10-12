// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidPreparation } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetValidPreparationResponseConfig extends ResponseConfig<ValidPreparation> {
		   validPreparationID: string;
		

		  constructor( validPreparationID: string, status: number = 200, body?: ValidPreparation) {
		    super();

		 this.validPreparationID = validPreparationID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetValidPreparation = (resCfg: MockGetValidPreparationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparations/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};