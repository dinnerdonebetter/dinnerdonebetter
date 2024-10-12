// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidPreparationVessel,
	QueryFilteredResult } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetValidPreparationVesselsByPreparationResponseConfig extends ResponseConfig<QueryFilteredResult<ValidPreparationVessel>> {
		   validPreparationID: string;
		

		  constructor( validPreparationID: string, status: number = 200, body: ValidPreparationVessel[] = []) {
		    super();

		 this.validPreparationID = validPreparationID;
		
		    this.status = status;
			if (this.body) {
			  this.body.data = body;
			}
		  }
}

export const mockGetValidPreparationVesselsByPreparations = (resCfg: MockGetValidPreparationVesselsByPreparationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparation_vessels/by_preparation/${resCfg.validPreparationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		
        if (resCfg.body && resCfg.filter) resCfg.body.limit = resCfg.filter.limit;
        if (resCfg.body && resCfg.filter) resCfg.body.page = resCfg.filter.page;
		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};