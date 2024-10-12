// GENERATED CODE, DO NOT EDIT MANUALLY

import type { Page, Route } from '@playwright/test';

import { ValidPreparationInstrument } from '@dinnerdonebetter/models'

import { assertClient, assertMethod, ResponseConfig } from './helpers';



export class MockArchiveValidPreparationInstrumentResponseConfig extends ResponseConfig<ValidPreparationInstrument> {
		   validPreparationVesselID: string;
		

		  constructor( validPreparationVesselID: string, status: number = 202, body?: ValidPreparationInstrument) {
		    super();

		 this.validPreparationVesselID = validPreparationVesselID;
		
		    this.status = status;
			if (this.body) {
			  this.body = body;
			}
		  }
}

export const mockArchiveValidPreparationInstrument = (resCfg: MockArchiveValidPreparationInstrumentResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/valid_preparation_instruments/${resCfg.validPreparationVesselID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('DELETE', route);
        assertClient(route);

		

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};