// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidPreparationInstrument } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockGetValidPreparationInstrumentResponseConfig extends ResponseConfig<ValidPreparationInstrument> {
		   validPreparationVesselID: string;
		

		  constructor( validPreparationVesselID: string, status: number = 200, body?: ValidPreparationInstrument) {
		    super();

		 this.validPreparationVesselID = validPreparationVesselID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockGetValidPreparationInstrument = (resCfg: MockGetValidPreparationInstrumentResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparation_instruments/${resCfg.validPreparationVesselID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};